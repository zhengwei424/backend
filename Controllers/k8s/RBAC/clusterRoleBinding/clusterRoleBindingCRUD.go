package clusterRoleBinding

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetClusterRoleBinding(c *gin.Context) {
	clusterRoleBinding := c.Query("clusterRoleBinding")
	clusterRoleBindingInfo := new(v1.ClusterRoleBinding)
	client := globalConfig.MyClient.Client
	clusterRoleBindingInfo, err := client.RbacV1().ClusterRoleBindings().Get(clusterRoleBinding, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		clusterRoleBindingInfo.Kind = "ClusterRoleBinding"
		clusterRoleBindingInfo.APIVersion = "rbac.authorization.k8s.io/v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": clusterRoleBindingInfo,
			"msg":  "ok",
		})
	}
}

func UpdateClusterRoleBinding(c *gin.Context) {
	clusterRoleBindingInfo := new(v1.ClusterRoleBinding)
	if err := c.BindJSON(clusterRoleBindingInfo); err == nil {
		fmt.Println(*clusterRoleBindingInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.RbacV1().ClusterRoleBindings().Update(clusterRoleBindingInfo)
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

func DeleteClusterRoleBinding(c *gin.Context) {
	clusterRoleBinding := c.Query("clusterRoleBinding")
	client := globalConfig.MyClient.Client
	err := client.RbacV1().ClusterRoleBindings().Delete(clusterRoleBinding, &metav1.DeleteOptions{})
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
