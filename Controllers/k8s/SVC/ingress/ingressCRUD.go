package ingress

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetIngress(c *gin.Context) {
	ns := c.Query("namespace")
	ingress := c.Query("ingress")
	ingressInfo := new(v1beta1.Ingress)

	client := globalConfig.MyClient.Client
	ingressInfo, err := client.ExtensionsV1beta1().Ingresses(ns).Get(ingress, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		ingressInfo.Kind = "Ingress"
		ingressInfo.APIVersion = "extensions/v1beta1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": ingressInfo,
			"msg":  "ok",
		})
	}
}

func UpdateIngress(c *gin.Context) {
	ingressInfo := new(v1beta1.Ingress)
	if err := c.BindJSON(ingressInfo); err == nil {
		fmt.Println(*ingressInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.ExtensionsV1beta1().Ingresses(ingressInfo.Namespace).Update(ingressInfo)
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

func DeleteIngress(c *gin.Context) {
	ns := c.Query("namespace")
	ingress := c.Query("ingress")

	client := globalConfig.MyClient.Client
	err := client.ExtensionsV1beta1().Ingresses(ns).Delete(ingress, &metav1.DeleteOptions{})
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
