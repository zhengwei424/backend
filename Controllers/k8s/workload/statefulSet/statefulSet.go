package statefulSet

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetStatefulSetsInfo(c *gin.Context) {
	var statefulSetsInfo = make([]map[string]interface{}, 0)
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
	statefulSets, err := client.AppsV1().StatefulSets(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, statefulSet := range statefulSets.Items {
		var statefulSetInfo = make(map[string]interface{}, 0)
		statefulSetInfo["name"] = statefulSet.Name
		statefulSetInfo["namespace"] = statefulSet.Namespace
		statefulSetInfo["labels"] = statefulSet.Labels
		statefulSetInfo["readyReplicas"] = statefulSet.Status.ReadyReplicas
		statefulSetInfo["replicas"] = statefulSet.Status.Replicas
		statefulSetInfo["creationTimestamp"] = tools.DeltaTime(statefulSet.CreationTimestamp.UTC(), time.Now())
		statefulSetsInfo = append(statefulSetsInfo, statefulSetInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": statefulSetsInfo,
		"msg":  "ok",
	})
}
