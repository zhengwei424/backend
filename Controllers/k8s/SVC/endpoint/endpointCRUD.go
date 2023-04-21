package endpoint

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetEndpoint(c *gin.Context) {
	ns := c.Query("namespace")
	endpoint := c.Query("endpoint")
	endpointInfo := new(v1.Endpoints)

	client := globalConfig.MyClient.Client
	endpointInfo, err := client.CoreV1().Endpoints(ns).Get(endpoint, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		endpointInfo.Kind = "endpoint"
		endpointInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": endpointInfo,
			"msg":  "ok",
		})
	}
}

func Updateendpoint(c *gin.Context) {
	endpointInfo := new(v1.Endpoints)
	if err := c.BindJSON(endpointInfo); err == nil {
		fmt.Println(*endpointInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().Endpoints(endpointInfo.Namespace).Update(endpointInfo)
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

func Deleteendpoint(c *gin.Context) {
	ns := c.Query("namespace")
	endpoint := c.Query("endpoint")

	client := globalConfig.MyClient.Client
	err := client.CoreV1().Endpoints(ns).Delete(endpoint, &metav1.DeleteOptions{})
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
