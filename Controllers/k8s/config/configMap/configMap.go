package configMap

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetConfigMapsInfo(c *gin.Context) {
	var configMapsInfo = make([]map[string]interface{}, 0)
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
	configMaps, err := client.CoreV1().ConfigMaps(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, configMap := range configMaps.Items {
		var configMapInfo = make(map[string]interface{}, 0)
		configMapInfo["name"] = configMap.Name
		configMapInfo["namespace"] = configMap.Namespace
		configMapInfo["labels"] = configMap.Labels
		configMapInfo["creationTimestamp"] = tools.DeltaTime(configMap.CreationTimestamp.UTC(), time.Now())
		configMapsInfo = append(configMapsInfo, configMapInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": configMapsInfo,
		"msg":  "ok",
	})
}
