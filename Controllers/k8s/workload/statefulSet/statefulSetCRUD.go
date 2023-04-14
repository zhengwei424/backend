package statefulSet

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetStatefulSet(c *gin.Context) {
	ns := c.Query("namespace")
	statefulSet := c.Query("statefulSet")
	statefulSetInfo := new(v1.StatefulSet)

	client := globalConfig.MyClient.Client
	statefulSetInfo, err := client.AppsV1().StatefulSets(ns).Get(statefulSet, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		statefulSetInfo.Kind = "StatefulSet"
		statefulSetInfo.APIVersion = "apps/v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": statefulSetInfo,
			"msg":  "ok",
		})
	}
}

func UpdateStatefulSet(c *gin.Context) {
	statefulSetInfo := new(v1.StatefulSet)
	if err := c.BindJSON(statefulSetInfo); err == nil {
		fmt.Println(*statefulSetInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.AppsV1().StatefulSets(statefulSetInfo.Namespace).Update(statefulSetInfo)
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

func DeleteStatefulSet(c *gin.Context) {
	ns := c.Query("namespace")
	statefulSet := c.Query("statefulSet")

	client := globalConfig.MyClient.Client
	err := client.AppsV1().StatefulSets(ns).Delete(statefulSet, &metav1.DeleteOptions{})
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
