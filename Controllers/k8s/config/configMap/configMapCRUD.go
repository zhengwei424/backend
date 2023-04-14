package configMap

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func CreateConfigMap(c *gin.Context) {
	configMapInfo := new(v1.ConfigMap)
	if err := c.BindJSON(configMapInfo); err == nil {
		fmt.Println(configMapInfo)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().ConfigMaps(configMapInfo.Namespace).Create(configMapInfo)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		//"data": configMapInfo,
		"msg": "ok",
	})
}

func GetConfigMap(c *gin.Context) {
	ns := c.Query("namespace")
	configMap := c.Query("configMap")
	configMapInfo := new(v1.ConfigMap)
	client := globalConfig.MyClient.Client
	configMapInfo, err := client.CoreV1().ConfigMaps(ns).Get(configMap, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		configMapInfo.Kind = "ConfigMap"
		configMapInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": configMapInfo,
			"msg":  "ok",
		})
	}
}

func UpdateConfigMap(c *gin.Context) {
	configMapInfo := new(v1.ConfigMap)
	if err := c.BindJSON(configMapInfo); err == nil {
		fmt.Println(*configMapInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().ConfigMaps(configMapInfo.Namespace).Update(configMapInfo)
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

func DeleteConfigMap(c *gin.Context) {
	ns := c.Query("namespace")
	configMap := c.Query("configMap")

	client := globalConfig.MyClient.Client
	err := client.CoreV1().ConfigMaps(ns).Delete(configMap, &metav1.DeleteOptions{})
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
