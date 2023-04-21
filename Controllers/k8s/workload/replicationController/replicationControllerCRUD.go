package replicationController

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetReplicationController(c *gin.Context) {
	ns := c.Query("namespace")
	replicationController := c.Query("replicationController")
	replicationControllerInfo := new(v1.ReplicationController)

	client := globalConfig.MyClient.Client
	replicationControllerInfo, err := client.CoreV1().ReplicationControllers(ns).Get(replicationController, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		replicationControllerInfo.Kind = "ReplicationController"
		replicationControllerInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": replicationControllerInfo,
			"msg":  "ok",
		})
	}
}

func UpdateReplicationController(c *gin.Context) {
	replicationControllerInfo := new(v1.ReplicationController)
	if err := c.BindJSON(replicationControllerInfo); err == nil {
		fmt.Println(*replicationControllerInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().ReplicationControllers(replicationControllerInfo.Namespace).Update(replicationControllerInfo)
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

func DeleteReplicationController(c *gin.Context) {
	ns := c.Query("namespace")
	replicationController := c.Query("replicationController")

	client := globalConfig.MyClient.Client
	err := client.CoreV1().ReplicationControllers(ns).Delete(replicationController, &metav1.DeleteOptions{})
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
