package role

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetRolesInfo(c *gin.Context) {
	var rolesInfo = make([]map[string]interface{}, 0)
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
	roles, err := client.RbacV1().Roles(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, role := range roles.Items {
		var roleInfo = make(map[string]interface{}, 0)
		roleInfo["name"] = role.Name
		roleInfo["namespace"] = role.Namespace
		roleInfo["labels"] = role.Labels
		roleInfo["rules"] = role.Rules
		roleInfo["creationTimestamp"] = tools.DeltaTime(role.CreationTimestamp.UTC(), time.Now())
		rolesInfo = append(rolesInfo, roleInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rolesInfo,
		"msg":  "ok",
	})
}
