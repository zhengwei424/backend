package roleBinding

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetRoleBindingsInfo(c *gin.Context) {
	var roleBindingsInfo = make([]map[string]interface{}, 0)
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
	roleBindings, err := client.RbacV1().RoleBindings(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, roleBinding := range roleBindings.Items {
		var roleBindingInfo = make(map[string]interface{}, 0)
		roleBindingInfo["name"] = roleBinding.Name
		roleBindingInfo["namespace"] = roleBinding.Namespace
		roleBindingInfo["labels"] = roleBinding.Labels
		roleBindingInfo["roleRef"] = roleBinding.RoleRef
		roleBindingInfo["subjects"] = roleBinding.Subjects
		roleBindingInfo["creationTimestamp"] = tools.DeltaTime(roleBinding.CreationTimestamp.UTC(), time.Now())
		roleBindingsInfo = append(roleBindingsInfo, roleBindingInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": roleBindingsInfo,
		"msg":  "ok",
	})
}
