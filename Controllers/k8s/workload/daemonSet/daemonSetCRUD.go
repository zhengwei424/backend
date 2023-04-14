package daemonSet

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetDaemonSet(c *gin.Context) {
	ns := c.Query("namespace")
	daemonSet := c.Query("daemonSet")
	daemonSetInfo := new(v1beta1.DaemonSet)

	client := globalConfig.MyClient.Client
	daemonSetInfo, err := client.ExtensionsV1beta1().DaemonSets(ns).Get(daemonSet, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		daemonSetInfo.Kind = "DaemonSet"
		daemonSetInfo.APIVersion = "extensions/v1beta1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": daemonSetInfo,
			"msg":  "ok",
		})
	}
}

func UpdateDaemonSet(c *gin.Context) {
	daemonSetInfo := new(v1beta1.DaemonSet)
	if err := c.BindJSON(daemonSetInfo); err == nil {
		fmt.Println(*daemonSetInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.ExtensionsV1beta1().DaemonSets(daemonSetInfo.Namespace).Update(daemonSetInfo)
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

func DeleteDaemonSet(c *gin.Context) {
	ns := c.Query("namespace")
	daemonSet := c.Query("daemonSet")

	client := globalConfig.MyClient.Client
	err := client.ExtensionsV1beta1().DaemonSets(ns).Delete(daemonSet, &metav1.DeleteOptions{})
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
