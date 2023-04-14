package secret

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func CreateSecret(c *gin.Context) {
	secretInfo := new(v1.Secret)
	if err := c.BindJSON(secretInfo); err == nil {
		fmt.Println(secretInfo)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().Secrets(secretInfo.Namespace).Create(secretInfo)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		//"data": secretInfo,
		"msg": "ok",
	})
}

func GetSecret(c *gin.Context) {
	ns := c.Query("namespace")
	secret := c.Query("secret")

	secretInfo := new(v1.Secret)
	client := globalConfig.MyClient.Client
	secretInfo, err := client.CoreV1().Secrets(ns).Get(secret, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		secretInfo.Kind = "Secret"
		secretInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": secretInfo,
			"msg":  "ok",
		})
	}
}

func UpdateSecret(c *gin.Context) {
	secretInfo := new(v1.Secret)
	if err := c.BindJSON(secretInfo); err == nil {
		fmt.Println(*secretInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().Secrets(secretInfo.Namespace).Update(secretInfo)
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

func DeleteSecret(c *gin.Context) {
	ns := c.Query("namespace")
	secret := c.Query("secret")

	client := globalConfig.MyClient.Client
	err := client.CoreV1().Secrets(ns).Delete(secret, &metav1.DeleteOptions{})
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
