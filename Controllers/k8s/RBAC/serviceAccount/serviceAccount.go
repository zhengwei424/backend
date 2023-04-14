package serviceAccount

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetServiceAccountsInfo(c *gin.Context) {
	var serviceAccountsInfo = make([]map[string]interface{}, 0)
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
	serviceAccounts, err := client.CoreV1().ServiceAccounts(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, serviceAccount := range serviceAccounts.Items {
		var serviceAccountInfo = make(map[string]interface{}, 0)
		serviceAccountInfo["name"] = serviceAccount.Name
		serviceAccountInfo["namespace"] = serviceAccount.Namespace
		serviceAccountInfo["labels"] = serviceAccount.Labels
		serviceAccountInfo["secrets"] = serviceAccount.Secrets
		serviceAccountInfo["creationTimestamp"] = tools.DeltaTime(serviceAccount.CreationTimestamp.UTC(), time.Now())
		serviceAccountsInfo = append(serviceAccountsInfo, serviceAccountInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": serviceAccountsInfo,
		"msg":  "ok",
	})
}
