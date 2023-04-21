package service

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetService(c *gin.Context) {
	ns := c.Query("namespace")
	service := c.Query("service")
	serviceInfo := new(v1.Service)

	client := globalConfig.MyClient.Client
	serviceInfo, err := client.CoreV1().Services(ns).Get(service, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		serviceInfo.Kind = "Service"
		serviceInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": serviceInfo,
			"msg":  "ok",
		})
	}
}

func UpdateService(c *gin.Context) {
	serviceInfo := new(v1.Service)
	if err := c.BindJSON(serviceInfo); err == nil {
		fmt.Println(*serviceInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().Services(serviceInfo.Namespace).Update(serviceInfo)
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

func DeleteService(c *gin.Context) {
	ns := c.Query("namespace")
	service := c.Query("service")

	client := globalConfig.MyClient.Client
	err := client.CoreV1().Services(ns).Delete(service, &metav1.DeleteOptions{})
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
