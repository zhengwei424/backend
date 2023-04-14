package replicaSet

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetReplicaSetsInfo(c *gin.Context) {
	var replicaSetsInfo = make([]map[string]interface{}, 0)
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
	replicaSets, err := client.AppsV1().ReplicaSets(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, replicaSet := range replicaSets.Items {
		var replicaSetInfo = make(map[string]interface{}, 0)
		replicaSetInfo["name"] = replicaSet.Name
		replicaSetInfo["namespace"] = replicaSet.Namespace
		replicaSetInfo["labels"] = replicaSet.Labels
		replicaSetInfo["creationTimestamp"] = tools.DeltaTime(replicaSet.CreationTimestamp.UTC(), time.Now())
		replicaSetsInfo = append(replicaSetsInfo, replicaSetInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": replicaSetsInfo,
		"msg":  "ok",
	})
}
