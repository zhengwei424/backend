package deprecated

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

/*
 message types:
 The message types are defined in RFC 6455, section 11.8.
 TextMessage denotes a text data message. The text message payload is
 interpreted as UTF-8 encoded text data.
 TextMessage = 1

 BinaryMessage denotes a binary data message.
 BinaryMessage = 2

 CloseMessage denotes a close control message. The optional message
 payload contains a numeric code and text. Use the FormatCloseMessage
 function to format a close message payload.
 CloseMessage = 8

 PingMessage denotes a ping control message. The optional message payload
 is UTF-8 encoded text.
 PingMessage = 9

 PongMessage denotes a pong control message. The optional message payload
 is UTF-8 encoded text.
 PongMessage = 10
*/

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有CORS跨域请求
		return true
	},
}

// websocket消息
type WsMessage struct {
	MessageType int
	Data        []byte
}

//封装websocket连接
type WsConnection struct {
	wsSocket *websocket.Conn // 底层webdocket
	inChan   chan *WsMessage // 读取队列(数据量大，就在通道中传递指针）
	outChan  chan *WsMessage // 发送队列(数据量大，就在通道中传递指针）

	mutex     sync.Mutex //互斥锁，避免重复关闭管道
	isClosed  bool
	closeChan chan byte // 关闭通知
}

// 读取协程
func (wsConn *WsConnection) wsReadLoop() {
	var (
		msgType int
		data    []byte
		msg     *WsMessage
		err     error
	)

	for {
		// 读一个message，客户端没有发送消息，服务端会在wsConn.wsSocket.ReadMessage()一直等待，类似与fmt.Scanln()效果，等待用户输入？？
		fmt.Println("读一个message")
		if msgType, data, err = wsConn.wsSocket.ReadMessage(); err != nil {
			// 当客户端没有消息传入时，会报错websocket: close 1005 (no status)
			fmt.Println("read message loop err:", err)
			goto ERROR
		}
		fmt.Println("读message完毕")
		msg = &WsMessage{
			msgType,
			data,
		}

		// 心跳检测
		// 获取客户端心跳并响应
		if string(msg.Data) == "{\"heartbeat\":\"ping\"}" {
			var pong string = "{\"heartbeat\":\"pong\"}"
			_ = wsConn.wsSocket.WriteMessage(msg.MessageType, []byte(pong))
			continue
		}

		fmt.Println("read loop msg:", string(data))
		// 放入请求队列
		select {
		// 将msg放入通道
		case wsConn.inChan <- msg:
			fmt.Print("msg进入通道")
		case <-wsConn.closeChan:
			goto CLOSED
		}
	}
ERROR:
	fmt.Println("执行WsClose")
	wsConn.WsClose()
CLOSED:
	// ????????????????
}

// 发送协程
func (wsConn *WsConnection) wsWriteLoop() {
	var (
		msg *WsMessage
		err error
	)
	for {
		select {
		// 取一个应答
		case msg = <-wsConn.outChan:
			// 写给websocket
			if err = wsConn.wsSocket.WriteMessage(msg.MessageType, msg.Data); err != nil {
				goto ERROR
			}
		case <-wsConn.closeChan:
			goto CLOSED
		}
	}
ERROR:
	wsConn.WsClose()
CLOSED:
	// ????????????????
}

// 向web端发送消息
func (wsConn *WsConnection) WsWrite(messageType int, data []byte) (err error) {
	select {
	case wsConn.outChan <- &WsMessage{messageType, data}:
	case <-wsConn.closeChan:
		err = errors.New("websocket closed")
	}
	return
}

func (wsConn *WsConnection) WsWriteJSON(v interface{}) (err error) {
	err = wsConn.wsSocket.WriteJSON(v)
	if err != nil {
		return err
	}
	return nil
}

// 读取web端消息
func (wsConn *WsConnection) WsRead() (msg *WsMessage, err error) {
	select {
	case msg = <-wsConn.inChan:
		return
	case <-wsConn.closeChan:
		err = errors.New("websocket closed")
	}
	return
}

//// pod日志专用，获取客户端查询参数
//func (wsConn *WsConnection) WsReadForPodLog(qryData *QueryData, block chan int) {
//	var msg *WsMessage
//	for {
//		fmt.Println("WsReadForPodLog is running")
//		select {
//		case msg = <-wsConn.inChan:
//			if err := json.Unmarshal(msg.Data, &qryData); err != nil {
//				fmt.Println("queryData err: ", err)
//			} else {
//				fmt.Println("in function qryData:", qryData)
//				block <- 1
//			}
//			//fmt.Println("in readpodlog function , msg.data:", string(msg.Data))
//		case <-wsConn.closeChan:
//			errors.New("websocket closed")
//		}
//	}
//
//}

func (wsConn *WsConnection) IsClosed() bool {
	return wsConn.isClosed
}

// 关闭连接
func (wsConn *WsConnection) WsClose() {
	// panic: use of closed network connection 使用一个已经关闭的网络连接
	err := wsConn.wsSocket.Close()
	if err != nil {
		panic(err)
	}

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}

// 并发安全API
func InitWebsocket(c *gin.Context) (wsConn *WsConnection, err error) {
	var (
		wsSocket *websocket.Conn
	)
	// 应答客户端，告知升级连接为websocket
	if wsSocket, err = wsUpgrader.Upgrade(c.Writer, c.Request, nil); err != nil {
		return
	}

	wsConn = &WsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *WsMessage, 1000),
		outChan:   make(chan *WsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
	}

	//// 读协程
	go wsConn.wsReadLoop()
	//// 写协程
	go wsConn.wsWriteLoop()

	return
}
