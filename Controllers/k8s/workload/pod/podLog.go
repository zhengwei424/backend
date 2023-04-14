package pod

import (
	"backend/globalConfig"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func HandleLogs(c *gin.Context) {
	var err error
	// 获取k8s自身资源管理客户端实例
	client := globalConfig.MyClient.Client

	namespace := c.Param("namespace")
	pod := c.Param("pod")
	container := c.Param("container")

	// 获取参考时间戳
	refTimestamp := c.Query("referenceTimestamp")
	if refTimestamp == "" {
		refTimestamp = NewestTimestamp
	}

	// 获取参考行
	refLineNum, err := strconv.Atoi(c.Query("referenceLineNum"))
	if err != nil {
		refLineNum = 0
	}

	// 是否显示以前terminated的容器的日志
	usePreviousLogs := c.Query("previous") == "true"
	//
	offsetFrom, err1 := strconv.Atoi(c.Query("offsetFrom"))
	//
	offsetTo, err2 := strconv.Atoi(c.Query("offsetTo"))
	//
	logFilePosition := c.Query("logFilePosition")
	//
	//follow := c.Query("follow")
	// 如果没有请求offset或者请求offset异常，logSelector使用默认配置
	logSelector := DefaultSelection
	if err1 == nil && err2 == nil {
		logSelector = &Selection{
			ReferencePoint: LogLineId{
				LogTimestamp: LogTimestamp(refTimestamp),
				LineNum:      refLineNum,
			},
			OffsetFrom:      offsetFrom,
			OffsetTo:        offsetTo,
			LogFilePosition: logFilePosition, // value is "beginning" or "end" ?
		}
	}

	var result = new(LogDetails)
	result, err = GetLogDetails(client, namespace, pod, container, logSelector, usePreviousLogs)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": result,
	})
}

func HandleLogFile(c *gin.Context) {
	var err error
	// 获取k8s自身资源管理客户端实例
	client := globalConfig.MyClient.Client

	namespace := c.Param("namespace")
	pod := c.Param("pod")
	container := c.Param("container")
	usePreviousLogs := c.Query("previous") == "true"

	logStream, err := GetLogFile(client, namespace, pod, container, usePreviousLogs)
	//var opt = &v1.PodLogOptions{
	//	Container: container,
	//	Follow:    false,
	//	Previous:  usePreviousLogs,
	//}
	//logStream, err := client.CoreV1().Pods(namespace).GetLogs(pod, opt).Stream()
	if err != nil {
		panic(err)
	}
	//var buf = new(bytes.Buffer)
	// logStream的内容会全部传到buf内，logStream内容为空，操作相当于mv
	//_, err = io.Copy(buf, logStream)
	//if err != nil {
	//	fmt.Println("err:", err.Error())
	//	panic(err)
	//}
	//
	//var p = make([]byte, 10)
	//_, _ = logStream.Read(p)
	//fmt.Println("logStream:", string(p))
	//fmt.Println("buff:", string(buf.Bytes()))
	//c.JSON(http.StatusOK, gin.H{
	//	"code": 0,
	//	"msg":  "ok",
	//})
	//c.DataFromReader(http.StatusOK, int64(buf.Len()), "text/plain", buf, map[string]string{})
	extraHeader := map[string]string{
		"filename": "logs-from-" + container + "-in-" + pod + ".log",
	}
	// contentLength和contentType就是响应头的内容，和logStream的传输无关，感觉只是反馈一个header信息而已
	c.DataFromReader(http.StatusOK, 10, "text/plain", logStream, extraHeader)
}
