package serviceAccount

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetServiceAccount(c *gin.Context) {
	ns := c.Query("namespace")
	serviceAccount := c.Query("serviceAccount")
	serviceAccountInfo := new(v1.ServiceAccount)
	client := globalConfig.MyClient.Client
	serviceAccountInfo, err := client.CoreV1().ServiceAccounts(ns).Get(serviceAccount, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		serviceAccountInfo.Kind = "ServiceAccount"
		serviceAccountInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": serviceAccountInfo,
			"msg":  "ok",
		})
	}
}

func UpdateServiceAccount(c *gin.Context) {
	serviceAccountInfo := new(v1.ServiceAccount)
	if err := c.BindJSON(serviceAccountInfo); err == nil {
		fmt.Println(*serviceAccountInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().ServiceAccounts(serviceAccountInfo.Namespace).Update(serviceAccountInfo)
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

func DeleteServiceAccount(c *gin.Context) {
	ns := c.Query("namespace")
	serviceAccount := c.Query("serviceAccount")

	client := globalConfig.MyClient.Client
	err := client.CoreV1().ServiceAccounts(ns).Delete(serviceAccount, &metav1.DeleteOptions{})
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
