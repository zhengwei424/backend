package pod

import (
	"backend/globalConfig"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
	"sync"
)

const END_OF_TRANSMISSION = "\u0004"

// TerminalMessage 是ShellController和TerminalSession之间的消息协议
// ---------------------------------------------------------------------
// fe表示前端？be表示后端？enpoint
// OP      DIRECTION  FIELD(S) USED  DESCRIPTION
// ---------------------------------------------------------------------
// bind    fe->be     SessionID      Id sent back from TerminalResponse
// stdin   fe->be     Data           Keystrokes/paste buffer
// resize  fe->be     Rows, Cols     New terminal size
// stdout  be->fe     Data           Output from the process
// toast   be->fe     Data           OOB message to be shown to the user
type TerminalMessage struct {
	Op, Data, SessionID string
	Rows, Cols          uint16
}

// PtyHandler 伪终端（pseudo-terminal）接口
type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

// TerminalSession 实现PtyHandler接口
type TerminalSession struct {
	id            string
	bound         chan error
	sockJSSession sockjs.Session
	sizeChan      chan remotecommand.TerminalSize
}

// Next 控制终端的resize事件
// 在进程运行过程中循环调用
func (t TerminalSession) Next() *remotecommand.TerminalSize {
	size := <-t.sizeChan
	if size.Height == 0 && size.Width == 0 {
		return nil
	}
	return &size
}

// Read 进程读取伪终端消息(stdin)
// 在进程运行过程中循环调用
func (t TerminalSession) Read(p []byte) (n int, err error) {
	// 接收终端的字符串
	m, err := t.sockJSSession.Recv()
	if err != nil {
		// 接收错误时，发送一个终止信号给进程，避免资源泄露
		return copy(p, END_OF_TRANSMISSION), err
	}

	var msg TerminalMessage
	if err := json.Unmarshal([]byte(m), &msg); err != nil {
		// json解析字符串到msg错误时，发送一个终止信号给进程
		return copy(p, END_OF_TRANSMISSION), err
	}

	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	default:
		return copy(p, END_OF_TRANSMISSION), fmt.Errorf("unkonwn message type '%s'", msg.Op)
	}
}

// Write 进程消息输出到终端(stdout)
// 在进程运行过程中循环调用
func (t TerminalSession) Write(p []byte) (n int, err error) {
	msg, err := json.Marshal(TerminalMessage{
		Op:   "stdout",
		Data: string(p),
	})
	if err != nil {
		return 0, err
	}

	if err = t.sockJSSession.Send(string(msg)); err != nil {
		return 0, err
	}

	return len(p), nil
}

// Toast 用于发送用户的任何OOB（out of browser ？）消息到终端(好像没啥用？？？？？？？？？？）
//func (t TerminalSession) Toast(p string) error {
//	msg, err := json.Marshal(TerminalMessage{
//		Op:   "toast",
//		Data: p,
//	})
//	if err != nil {
//		return err
//	}
//
//	if err := t.sockJSSession.Send(string(msg)); err != nil {
//		return err
//	}
//
//	return nil
//}

// SessionMap 存储所有TerminalSession
// Lock 避免并发冲突
type SessionMap struct {
	Sessions map[string]TerminalSession
	Lock     sync.RWMutex
}

// Get 通过sessionID获取对应的TerminalSession
func (sm *SessionMap) Get(sessionID string) TerminalSession {
	sm.Lock.RLock()
	defer sm.Lock.RUnlock()
	return sm.Sessions[sessionID]
}

// Set 将session添加到sessionMap中
func (sm *SessionMap) Set(sessionID string, session TerminalSession) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	sm.Sessions[sessionID] = session
}

// Close 关闭指定的session
func (sm *SessionMap) Close(sessionID string, status uint32, reason string) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	// 通过sessionID获取session
	ses := sm.Sessions[sessionID]
	// 调用Session接口的Close方法
	err := ses.sockJSSession.Close(status, reason)
	if err != nil {
		log.Print(err)
	}
	// 关闭sizeChan
	close(ses.sizeChan)
	// 从sessionMap中移除该session
	delete(sm.Sessions, sessionID)
}

// 初始化sessionMap
var terminalSessions = SessionMap{Sessions: make(map[string]TerminalSession)}

// 被net/http调用，用于处理任何请求/api/sockjs的连接，绑定sessionMap中存在的session
func handleTerminalSession(session sockjs.Session) {
	var (
		buf             string
		err             error
		msg             TerminalMessage
		terminalSession TerminalSession
	)

	// 读取session内容，返回string
	if buf, err = session.Recv(); err != nil {
		log.Printf("handleTerminalSession: can`t Recv: %v", err)
		return
	}

	// 解析buf为msg
	if err = json.Unmarshal([]byte(buf), &msg); err != nil {
		log.Printf("handleTerminalSession: can`t UnMarshal (%v): %s", err, buf)
		return
	}

	// 解析msg，判断是否时bind操作
	if msg.Op != "bind" {
		log.Printf("handleTerminalSession: expected 'bind' message, got: %s", buf)
		return
	}

	// 确认sessionID为非空
	if terminalSession = terminalSessions.Get(msg.SessionID); terminalSession.id == "" {
		log.Printf("handleTerminalSession: can`t find session '%s'", msg.SessionID)
		return
	}

	terminalSession.sockJSSession = session
	terminalSessions.Set(msg.SessionID, terminalSession)
	terminalSession.bound <- nil
}

func CreateAttachHandler(path string) http.Handler {
	return sockjs.NewHandler(path, sockjs.DefaultOptions, handleTerminalSession)
}

func startProcess(k8sClient kubernetes.Interface, cfg *rest.Config, c *gin.Context, cmd []string, ptyHandler PtyHandler) error {
	namespace := c.Param("namespace")
	podName := c.Param("pod")
	containerName := c.Param("container")

	req := k8sClient.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(namespace).SubResource("exec")

	req.VersionedParams(&v1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	//fmt.Println("apipath:", cfg.APIPath)
	//fmt.Println("ca:", string(cfg.CAData))
	//fmt.Println("cert:", string(cfg.CertData))
	//fmt.Println("key:", string(cfg.KeyData))
	exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		fmt.Println("err1:", err)
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: ptyHandler,
		Tty:               true,
	})
	if err != nil {
		return err
	}
	return nil
}

// 生成一个随机的sessionID
func genTerminalSessionId() (string, error) {
	bytes := make([]byte, 16)
	// 生成一个长度为16的随机的bytes
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// hex.EncodedLen(len(bytes))获取bytes编码后的长度
	id := make([]byte, hex.EncodedLen(len(bytes)))
	hex.Encode(id, bytes)
	return string(id), nil
}

// 判断shell是否在validShells数组内
func isValidShell(validShells []string, shell string) bool {
	for _, validShell := range validShells {
		if validShell == shell {
			return true
		}
	}
	return false
}

func WaitForTerminal(k8sClient kubernetes.Interface, cfg *rest.Config, c *gin.Context, sessionID string) {
	shell := c.Query("shell")

	// 复制原始上下文，用于参数传递
	cCp := c.Copy()

	select {
	case <-terminalSessions.Get(sessionID).bound:
		close(terminalSessions.Get(sessionID).bound)

		var err error
		validShells := []string{"bash", "sh", "powershell", "cmd"}

		if isValidShell(validShells, shell) {
			cmd := []string{shell}
			err = startProcess(k8sClient, cfg, cCp, cmd, terminalSessions.Get(sessionID))
		} else {
			// 没有传入shell时，就遍历validShells，哪个能用用哪个
			for _, testShell := range validShells {
				cmd := []string{testShell}
				if err = startProcess(k8sClient, cfg, cCp, cmd, terminalSessions.Get(sessionID)); err == nil {
					break
				}
			}
		}

		if err != nil {
			terminalSessions.Close(sessionID, 2, err.Error())
			return
		}

		terminalSessions.Close(sessionID, 2, "Process exited")
	}
}

func HandleExecShell(c *gin.Context) {
	// 生成sessionID
	sessionID, err := genTerminalSessionId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	// 获取k8sClient
	client := globalConfig.MyClient.Client

	// 获取kubeconfig配置？
	cfg := globalConfig.MyCfg
	//fmt.Println("apipath:", cfg.APIPath)
	//fmt.Println("ca:", string(cfg.CAData))
	//fmt.Println("cert:", string(cfg.CertData))
	//fmt.Println("key:", string(cfg.KeyData))

	terminalSessions.Set(sessionID, TerminalSession{
		id:       sessionID,
		bound:    make(chan error),
		sizeChan: make(chan remotecommand.TerminalSize),
	})

	// 复制原始上下文，用于参数传递!!!!!!!
	cCp := c.Copy()

	go WaitForTerminal(client, cfg, cCp, sessionID)
	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"sessionID": sessionID,
	})
}
