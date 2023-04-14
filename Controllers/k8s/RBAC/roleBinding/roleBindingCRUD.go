package roleBinding

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func CreateRoleBinding(c *gin.Context) {
	roleBindingInfo := new(v1.RoleBinding)
	if err := c.BindJSON(roleBindingInfo); err == nil {
		fmt.Println(roleBindingInfo)
	}

	client := globalConfig.MyClient.Client
	_, err := client.RbacV1().RoleBindings(roleBindingInfo.Namespace).Create(roleBindingInfo)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		//"data": roleBindingInfo,
		"msg": "ok",
	})
}

func GetRoleBinding(c *gin.Context) {
	ns := c.Query("namespace")
	roleBinding := c.Query("roleBinding")
	roleBindingInfo := new(v1.RoleBinding)
	client := globalConfig.MyClient.Client
	roleBindingInfo, err := client.RbacV1().RoleBindings(ns).Get(roleBinding, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		roleBindingInfo.Kind = "RoleBinding"
		roleBindingInfo.APIVersion = "rbac.authorization.k8s.io/v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": roleBindingInfo,
			"msg":  "ok",
		})
	}
}

func UpdateRoleBinding(c *gin.Context) {
	roleBindingInfo := new(v1.RoleBinding)
	if err := c.BindJSON(roleBindingInfo); err == nil {
		fmt.Println(*roleBindingInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.RbacV1().RoleBindings(roleBindingInfo.Namespace).Update(roleBindingInfo)
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

func DeleteRoleBinding(c *gin.Context) {
	ns := c.Query("namespace")
	roleBinding := c.Query("roleBinding")

	client := globalConfig.MyClient.Client
	err := client.RbacV1().RoleBindings(ns).Delete(roleBinding, &metav1.DeleteOptions{})
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
