package pod

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func CreatePod(c *gin.Context) {
	podInfo := new(v1.Pod)
	if err := c.BindJSON(podInfo); err == nil {
		fmt.Println(podInfo.Status)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().Pods(podInfo.Namespace).Create(podInfo)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		//"data": podInfo,
		"msg": "ok",
	})
}

func GetPod(c *gin.Context) {
	ns := c.Query("namespace")
	pod := c.Query("pod")
	podInfo := new(v1.Pod)

	client := globalConfig.MyClient.Client
	podInfo, err := client.CoreV1().Pods(ns).Get(pod, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		podInfo.Kind = "Pod"
		podInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": podInfo,
			"msg":  "ok",
		})
	}
}

func UpdatePod(c *gin.Context) {
	podInfo := new(v1.Pod)
	if err := c.BindJSON(podInfo); err == nil {
		fmt.Println(*podInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().Pods(podInfo.Namespace).Update(podInfo)
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

func DeletePod(c *gin.Context) {
	ns := c.Query("namespace")
	pod := c.Query("pod")

	client := globalConfig.MyClient.Client
	err := client.CoreV1().Pods(ns).Delete(pod, &metav1.DeleteOptions{})
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
