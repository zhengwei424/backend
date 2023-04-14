package clusterRole

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetClusterRolesInfo(c *gin.Context) {
	var clusterRolesInfo = make([]map[string]interface{}, 0)
	// 获取k8s自身资源管理客户端实例
	client := globalConfig.MyClient.Client

	opts := v1.ListOptions{}
	clusterRoles, err := client.RbacV1().ClusterRoles().List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, clusterRole := range clusterRoles.Items {
		var clusterRoleInfo = make(map[string]interface{}, 0)
		clusterRoleInfo["name"] = clusterRole.Name
		clusterRoleInfo["labels"] = clusterRole.Labels
		clusterRoleInfo["rules"] = clusterRole.Rules
		clusterRoleInfo["creationTimestamp"] = tools.DeltaTime(clusterRole.CreationTimestamp.UTC(), time.Now())
		clusterRolesInfo = append(clusterRolesInfo, clusterRoleInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": clusterRolesInfo,
		"msg":  "ok",
	})
}
