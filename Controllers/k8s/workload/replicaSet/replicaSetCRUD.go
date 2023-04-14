package replicaSet

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetReplicaSet(c *gin.Context) {
	ns := c.Query("namespace")
	replicaSet := c.Query("replicaSet")
	replicaSetInfo := new(v1beta1.ReplicaSet)

	client := globalConfig.MyClient.Client
	replicaSetInfo, err := client.ExtensionsV1beta1().ReplicaSets(ns).Get(replicaSet, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		replicaSetInfo.Kind = "ReplicaSet"
		replicaSetInfo.APIVersion = "extensions/v1beta1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": replicaSetInfo,
			"msg":  "ok",
		})
	}
}

func UpdateReplicaSet(c *gin.Context) {
	replicaSetInfo := new(v1beta1.ReplicaSet)
	if err := c.BindJSON(replicaSetInfo); err == nil {
		fmt.Println(*replicaSetInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.ExtensionsV1beta1().ReplicaSets(replicaSetInfo.Namespace).Update(replicaSetInfo)
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

func DeleteReplicaSet(c *gin.Context) {
	ns := c.Query("namespace")
	replicaSet := c.Query("replicaSet")

	client := globalConfig.MyClient.Client
	err := client.ExtensionsV1beta1().ReplicaSets(ns).Delete(replicaSet, &metav1.DeleteOptions{})
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
