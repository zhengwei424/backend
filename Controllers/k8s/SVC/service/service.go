package service

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetServicesInfo(c *gin.Context) {
	var servicesInfo = make([]map[string]interface{}, 0)
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
	services, err := client.CoreV1().Services(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, service := range services.Items {
		var serviceInfo = make(map[string]interface{}, 0)
		serviceInfo["name"] = service.Name
		serviceInfo["namespace"] = service.Namespace
		serviceInfo["labels"] = service.Labels
		serviceInfo["creationTimestamp"] = tools.DeltaTime(service.CreationTimestamp.UTC(), time.Now())
		servicesInfo = append(servicesInfo, serviceInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": servicesInfo,
		"msg":  "ok",
	})
}
