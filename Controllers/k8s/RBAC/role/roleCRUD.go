package role

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetRole(c *gin.Context) {
	ns := c.Query("namespace")
	role := c.Query("role")
	roleInfo := new(v1.Role)
	client := globalConfig.MyClient.Client
	roleInfo, err := client.RbacV1().Roles(ns).Get(role, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		roleInfo.Kind = "Role"
		roleInfo.APIVersion = "rbac.authorization.k8s.io/v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": roleInfo,
			"msg":  "ok",
		})
	}
}

func UpdateRole(c *gin.Context) {
	roleInfo := new(v1.Role)
	if err := c.BindJSON(roleInfo); err == nil {
		fmt.Println(*roleInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.RbacV1().Roles(roleInfo.Namespace).Update(roleInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "ok",
		})
	}
}

func DeleteRole(c *gin.Context) {
	ns := c.Query("namespace")
	role := c.Query("role")

	client := globalConfig.MyClient.Client
	err := client.RbacV1().Roles(ns).Delete(role, &metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "ok",
		})
	}
}
