package clusterRoleBinding

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
	"time"
)

func GetClusterRoleBindingsInfo(c *gin.Context) {
	var clusterRoleBindingsInfo = make([]map[string]interface{}, 0)
	// 获取k8s自身资源管理客户端实例
	client := globalConfig.MyClient.Client

	opts := v1.ListOptions{}
	clusterRoleBindings, err := client.RbacV1().ClusterRoleBindings().List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		var clusterRoleBindingInfo = make(map[string]interface{}, 0)
		clusterRoleBindingInfo["name"] = clusterRoleBinding.Name
		clusterRoleBindingInfo["labels"] = clusterRoleBinding.Labels
		bindings := make([]string, 0)
		for _, v := range clusterRoleBinding.Subjects {
			bindings = append(bindings, v.Name)
		}
		clusterRoleBindingInfo["bindings"] = strings.Join(bindings, ",")
		clusterRoleBindingInfo["creationTimestamp"] = tools.DeltaTime(clusterRoleBinding.CreationTimestamp.UTC(), time.Now())
		clusterRoleBindingsInfo = append(clusterRoleBindingsInfo, clusterRoleBindingInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": clusterRoleBindingsInfo,
		"msg":  "ok",
	})
}
