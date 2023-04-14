package event

import (
	"backend/globalConfig"
	"backend/tools"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetEventsInfo(c *gin.Context) {
	var eventsInfo = make([]map[string]interface{}, 0)
	var qry, ns string
	// 获取k8s自身资源管理客户端实例
	client := globalConfig.MyClient.Client

	// 获取url中的namespace请求参数
	qry = c.Query("namespace")
	if qry == "all" {
		ns = ""
	} else {
		ns = qry
	}

	opts := v1.ListOptions{}
	events, err := client.CoreV1().Events(ns).List(opts)
	if err != nil {
		panic(err)
	}
	for _, event := range events.Items {
		var eventInfo = make(map[string]interface{}, 0)
		eventInfo["name"] = event.Name
		eventInfo["namespace"] = event.Namespace
		eventInfo["type"] = event.Type
		eventInfo["reason"] = event.Reason
		eventInfo["message"] = event.Message
		eventInfo["involveObject_kind"] = event.InvolvedObject.Kind
		eventInfo["involveObject_name"] = event.InvolvedObject.Name
		eventInfo["involveObject_namespace"] = event.InvolvedObject.Namespace
		eventInfo["source_component"] = event.Source.Component
		eventInfo["source_host"] = event.Source.Host
		eventInfo["first_timestamp"] = event.FirstTimestamp.UTC().String()
		eventInfo["last_timestamp"] = event.LastTimestamp.UTC().String()
		eventInfo["count"] = event.Count
		// 当前时间与lastTimestamp的时间差
		eventInfo["age"] = tools.DeltaTime(event.LastTimestamp.UTC(), time.Now())
		// 当前时间与firstTimestamp的时间差
		eventInfo["fullAge"] = tools.DeltaTime(event.FirstTimestamp.UTC(), time.Now())
		eventsInfo = append(eventsInfo, eventInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": eventsInfo,
		"msg":  "ok",
	})
}
