package deployment

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetDeploymentsInfo(c *gin.Context) {
	var deploymentsInfo = make([]map[string]interface{}, 0)
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
	deployments, err := client.AppsV1().Deployments(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, deployment := range deployments.Items {
		var deploymentInfo = make(map[string]interface{}, 0)
		deploymentInfo["name"] = deployment.Name
		deploymentInfo["namespace"] = deployment.Namespace
		deploymentInfo["labels"] = deployment.Labels
		deploymentInfo["creationTimestamp"] = tools.DeltaTime(deployment.CreationTimestamp.UTC(), time.Now())
		deploymentsInfo = append(deploymentsInfo, deploymentInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": deploymentsInfo,
		"msg":  "ok",
	})
}
