package daemonSet

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetDaemonSetsInfo(c *gin.Context) {
	var daemonSetsInfo = make([]map[string]interface{}, 0)
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
	daemonSets, err := client.AppsV1().DaemonSets(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, daemonSet := range daemonSets.Items {
		var daemonSetInfo = make(map[string]interface{}, 0)
		daemonSetInfo["name"] = daemonSet.Name
		daemonSetInfo["namespace"] = daemonSet.Namespace
		daemonSetInfo["labels"] = daemonSet.Labels
		daemonSetInfo["creationTimestamp"] = tools.DeltaTime(daemonSet.CreationTimestamp.UTC(), time.Now())
		daemonSetsInfo = append(daemonSetsInfo, daemonSetInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": daemonSetsInfo,
		"msg":  "ok",
	})
}
