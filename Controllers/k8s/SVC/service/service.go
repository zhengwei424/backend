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
		serviceInfo["type"] = service.Spec.Type
		serviceInfo["clusterIP"] = service.Spec.ClusterIP
		ports := make([]string, 0)
		for _, item := range service.Spec.Ports {
			var port string
			if item.NodePort != 0 {
				port = fmt.Sprintf("%d:%d/%s", item.Port, item.NodePort, item.Protocol)
			} else {
				if item.Port == item.TargetPort.IntVal {
					port = fmt.Sprintf("%d/%s", item.Port, item.Protocol)
				} else {
					port = fmt.Sprintf("%d:%d/%s", item.Port, item.TargetPort.IntVal, item.Protocol)
				}
			}
			ports = append(ports, port)
		}
		serviceInfo["ports"] = ports
		serviceInfo["externalIP"] = service.Spec.ExternalIPs
		serviceInfo["status"] = service.Spec.Type != "LoadBalancer" || len(service.Spec.ExternalIPs) > 0
		serviceInfo["selector"] = service.Spec.Selector
		serviceInfo["creationTimestamp"] = tools.DeltaTime(service.CreationTimestamp.UTC(), time.Now())
		servicesInfo = append(servicesInfo, serviceInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": servicesInfo,
		"msg":  "ok",
	})
}
