package secret

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetSecretsInfo(c *gin.Context) {
	var secretsInfo = make([]map[string]interface{}, 0)
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
	secrets, err := client.CoreV1().Secrets(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, secret := range secrets.Items {
		var secretInfo = make(map[string]interface{}, 0)
		secretInfo["name"] = secret.Name
		secretInfo["namespace"] = secret.Namespace
		secretInfo["labels"] = secret.Labels
		secretInfo["creationTimestamp"] = tools.DeltaTime(secret.CreationTimestamp.UTC(), time.Now())
		secretsInfo = append(secretsInfo, secretInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": secretsInfo,
		"msg":  "ok",
	})
}
