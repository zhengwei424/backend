package event

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetEvent(c *gin.Context) {
	ns := c.Query("namespace")
	event := c.Query("event")

	eventInfo := new(v1.Event)
	client := globalConfig.MyClient.Client
	eventInfo, err := client.CoreV1().Events(ns).Get(event, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		eventInfo.Kind = "Event"
		eventInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": eventInfo,
			"msg":  "ok",
		})
	}
}

func UpdateEvent(c *gin.Context) {
	eventInfo := new(v1.Event)
	if err := c.BindJSON(eventInfo); err == nil {
		fmt.Println(*eventInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().Events(eventInfo.Namespace).Update(eventInfo)
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

func DeleteEvent(c *gin.Context) {
	ns := c.Query("namespace")
	event := c.Query("event")

	client := globalConfig.MyClient.Client
	err := client.CoreV1().Events(ns).Delete(event, &metav1.DeleteOptions{})
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
