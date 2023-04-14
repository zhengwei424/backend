package clusterRole

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetClusterRole(c *gin.Context) {
	clusterRole := c.Query("clusterRole")
	clusterRoleInfo := new(v1.ClusterRole)
	client := globalConfig.MyClient.Client
	clusterRoleInfo, err := client.RbacV1().ClusterRoles().Get(clusterRole, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		clusterRoleInfo.Kind = "ClusterRole"
		clusterRoleInfo.APIVersion = "rbac.authorization.k8s.io/v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": clusterRoleInfo,
			"msg":  "ok",
		})
	}
}

func UpdateClusterRole(c *gin.Context) {
	clusterRoleInfo := new(v1.ClusterRole)
	if err := c.BindJSON(clusterRoleInfo); err == nil {
		fmt.Println(*clusterRoleInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.RbacV1().ClusterRoles().Update(clusterRoleInfo)
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

func DeleteClusterRole(c *gin.Context) {
	clusterRole := c.Query("clusterRole")

	client := globalConfig.MyClient.Client
	err := client.RbacV1().ClusterRoles().Delete(clusterRole, &metav1.DeleteOptions{})
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
