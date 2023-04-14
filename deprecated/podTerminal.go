package deprecated

import (
	"backend/globalConfig"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	clientset *kubernetes.Clientset
)

// ssh流式处理器
type streamHandler struct {
	wsConn      *WsConnection
	resizeEvent chan remotecommand.TerminalSize
}

// web终端发来的包
type xtermMessage struct {
	MsgType string `json:"type"`  // 类型:resize客户端调整终端, input客户端输入
	Input   string `json:"input"` // msgtype=input情况下使用
	Rows    uint16 `json:"rows"`  // msgtype=resize情况下使用
	Cols    uint16 `json:"cols"`  // msgtype=resize情况下使用
}

// Next executor回调获取web是否resize
func (handler *streamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <-handler.resizeEvent
	size = &ret
	return
}

// executor回调读取web端的输入
func (handler *streamHandler) Read(p []byte) (size int, err error) {
	var (
		msg      *WsMessage
		xtermMsg xtermMessage
	)

	// 读web发来的输入
	if msg, err = handler.wsConn.WsRead(); err != nil {
		return
	}

	// 解析客户端请求
	if err = json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		return
	}
	//web ssh调整了终端大小
	if xtermMsg.MsgType == "resize" {
		// 放到channel里，等remotecommand executor调用我们的Next取走
		handler.resizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	} else if xtermMsg.MsgType == "input" { // web ssh终端输入了字符
		// copy到p数组中
		size = len(xtermMsg.Input)
		copy(p, xtermMsg.Input)
	}
	return
}

// executor回调向web端输出
func (handler *streamHandler) Write(p []byte) (size int, err error) {
	var (
		copyData []byte
	)
	// 产生副本
	copyData = make([]byte, len(p))
	copy(copyData, p)
	size = len(p)
	err = handler.wsConn.WsWrite(websocket.TextMessage, copyData)
	return
}

func WsHandler(c *gin.Context) {
	var (
		wsConn    *WsConnection
		sshReq    *rest.Request
		pod       string
		namespace string
		container string
		shell     string
		executor  remotecommand.Executor
		handler   *streamHandler
		err       error
	)

	// 解析GET参数
	if err = c.Request.ParseForm(); err != nil {
		return
	}
	// 使用Form属性获取GET方法中的query参数 "/workload/exec?namespace=default&pod=nginx&container=nginx"
	namespace = c.Request.Form.Get("namespace")
	pod = c.Request.Form.Get("pod")
	container = c.Request.Form.Get("container")
	shell = c.Request.Form.Get("shell")

	// 获取pods

	// 获取k8s restclient以及kubeconfig配置
	client := globalConfig.MyClient.Client
	cfg := globalConfig.MyCfg

	// URL长相:
	// https://192.168.10.150:6443/api/v1/namespaces/default/pods/nginx/exec?command=bash&container=nginx&stderr=true&stdin=true&stdout=true&tty=true
	sshReq = client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: container,
			Command:   []string{shell},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	// 得到websocket长连接
	if wsConn, err = InitWebsocket(c); err != nil {
		return
	}

	// 创建到容器的连接
	if executor, err = remotecommand.NewSPDYExecutor(cfg, "POST", sshReq.URL()); err != nil {
		fmt.Println("remotecommand.NewSPDYExecutor error: ", err)
		goto END
	}

	// 配置与容器之间的数据流处理回调
	handler = &streamHandler{wsConn: wsConn, resizeEvent: make(chan remotecommand.TerminalSize)}
	if err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		goto END
	}
	return

END:
	wsConn.WsClose()
}
