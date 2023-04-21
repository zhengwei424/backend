package endpoint

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetEndpointsInfo(c *gin.Context) {
	var endpointsInfo = make([]map[string]interface{}, 0)
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
	endpoints, err := client.CoreV1().Endpoints(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, endpoint := range endpoints.Items {
		var endpointInfo = make(map[string]interface{}, 0)
		endpointInfo["name"] = endpoint.Name
		endpointInfo["namespace"] = endpoint.Namespace
		endpointInfo["labels"] = endpoint.Labels
		endpointInfo["creationTimestamp"] = tools.DeltaTime(endpoint.CreationTimestamp.UTC(), time.Now())
		endpointsInfo = append(endpointsInfo, endpointInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": endpointsInfo,
		"msg":  "ok",
	})
}
