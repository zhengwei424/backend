package ingress

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetIngressesInfo(c *gin.Context) {
	var ingressesInfo = make([]map[string]interface{}, 0)
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
	ingresses, err := client.ExtensionsV1beta1().Ingresses(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, ingress := range ingresses.Items {
		var ingressInfo = make(map[string]interface{}, 0)
		ingressInfo["name"] = ingress.Name
		ingressInfo["namespace"] = ingress.Namespace
		ingressInfo["labels"] = ingress.Labels
		ingressInfo["creationTimestamp"] = tools.DeltaTime(ingress.CreationTimestamp.UTC(), time.Now())
		ingressesInfo = append(ingressesInfo, ingressInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": ingressesInfo,
		"msg":  "ok",
	})
}
