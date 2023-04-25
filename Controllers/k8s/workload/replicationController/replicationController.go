package replicationController

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetReplicationControllersInfo(c *gin.Context) {
	var replicationControllersInfo = make([]map[string]interface{}, 0)
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
	replicationControllers, err := client.CoreV1().ReplicationControllers(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, replicationController := range replicationControllers.Items {
		var replicationControllerInfo = make(map[string]interface{}, 0)
		replicationControllerInfo["name"] = replicationController.Name
		replicationControllerInfo["namespace"] = replicationController.Namespace
		replicationControllerInfo["labels"] = replicationController.Labels
		replicationControllerInfo["readyReplicas"] = replicationController.Status.ReadyReplicas
		replicationControllerInfo["replicas"] = replicationController.Status.Replicas
		replicationControllerInfo["creationTimestamp"] = tools.DeltaTime(replicationController.CreationTimestamp.UTC(), time.Now())
		replicationControllersInfo = append(replicationControllersInfo, replicationControllerInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": replicationControllersInfo,
		"msg":  "ok",
	})
}
