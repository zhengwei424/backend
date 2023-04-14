package deprecated

import (
	"backend/globalConfig"
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	v1 "k8s.io/api/core/v1"
	"strconv"
	"time"
)

// QueryData 接收来自客户端的查询参数
type QueryData struct {
	Namespace   string `json:"namespace"`
	Pod         string `json:"pod"`
	Container   string `json:"container"`
	IsFollow    bool   `json:"isFollow"`
	TailLines   int64  `json:"tailLines"`
	PageSize    int64  `json:"PageSize"`
	CurrentPage int64  `json:"currentPage"`
}

func GetContainerLog(c *gin.Context) {
	var (
		client  = globalConfig.MyClient.Client
		wsConn  *WsConnection
		qryData *QueryData
		req     *io.ReadCloser
		err     error
	)
	// new初始化指针，避免空指针
	wsConn = new(WsConnection)
	qryData = new(QueryData)
	req = new(io.ReadCloser)
	//初始化websocket
	if wsConn, err = InitWebsocket(c); err != nil {
		panic(err)
	}

	// 从url获取请求参数
	//test := c.Request.URL.Query()
	//fmt.Println("test:", test)
	qryData.Namespace = c.Query("namespace")
	qryData.Pod = c.Query("pod")
	qryData.Container = c.Query("container")
	if qryData.PageSize, err = strconv.ParseInt(c.Query("pageSize"), 10, 0); err != nil {
		panic(err)
	}
	if qryData.CurrentPage, err = strconv.ParseInt(c.Query("currentPage"), 10, 0); err != nil {
		panic(err)
	}
	if qryData.Namespace == "" || qryData.Pod == "" || qryData.Container == "" {
		panic("namespace || pod || container can not be empty")
	}
	if qryData.IsFollow, err = strconv.ParseBool(c.Query("isFollow")); err != nil {
		fmt.Println("follow is not bool type")
		panic(err)
	}
	if qryData.TailLines, err = strconv.ParseInt(c.Query("tailLines"), 10, 64); err != nil {
		fmt.Println("tailLines is not int64 type")
	}

	// 获取日志流数据
	opts := &v1.PodLogOptions{
		Container: qryData.Container,
		Follow:    qryData.IsFollow,
		TailLines: &qryData.TailLines,
	}
	*req, err = client.CoreV1().Pods(qryData.Namespace).GetLogs(qryData.Pod, opts).Stream()
	if err != nil {
		err = (*req).Close()
	}
	r := bufio.NewReader(*req)

	/* 日志follow时，不分页 */
	for qryData.IsFollow {
		line, err := r.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		}
		err = wsConn.WsWrite(websocket.TextMessage, line)
		if err != nil {
			fmt.Println(err)
		}
	}

	/* 日志不follow时，分页 */
	// 单页日志
	var pagelines = make([][]byte, 0)
	// 日志总行数
	var totalLines int64 = 0
	// 满足pageSize的页数
	var fullPageNum int64 = 0
	// 存放分页数据
	var pages = make([][][]byte, 0)
	for !qryData.IsFollow {
		// 单行日志
		line, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// 数据已经读完了
				// 返回分页数据
				if fullPageNum < 1 || qryData.CurrentPage > fullPageNum {
					pages = append(pages, pagelines)
				}

				err = wsConn.WsWriteJSON(gin.H{"totalLines": totalLines})
				for _, l := range pages[qryData.CurrentPage-1] {
					err = wsConn.WsWrite(websocket.TextMessage, l)
				}
				if err != nil {
					fmt.Println("write to web auth error:", err)
				}
				// ????????????????????
				time.Sleep(24 * time.Hour)
			}
		} else {
			totalLines += 1
			pagelines = append(pagelines, line)
			if int64(len(pagelines)) == qryData.PageSize {
				fullPageNum += 1
				pages = append(pages, pagelines)
				pagelines = nil
			}
		}
	}

}

// [GIN-debug] [WARNING] Headers were already written. Wanted to override status code 200 with 500
