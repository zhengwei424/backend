package deployment

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetDeployment(c *gin.Context) {
	ns := c.Query("namespace")
	deployment := c.Query("deployment")
	deploymentInfo := new(v1beta1.Deployment)

	client := globalConfig.MyClient.Client
	deploymentInfo, err := client.ExtensionsV1beta1().Deployments(ns).Get(deployment, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		deploymentInfo.Kind = "Deployment"
		deploymentInfo.APIVersion = "extensions/v1beta1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": deploymentInfo,
			"msg":  "ok",
		})
	}
}

func UpdateDeployment(c *gin.Context) {
	deploymentInfo := new(v1beta1.Deployment)
	if err := c.BindJSON(deploymentInfo); err == nil {
		fmt.Println(*deploymentInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.ExtensionsV1beta1().Deployments(deploymentInfo.Namespace).Update(deploymentInfo)
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

func DeleteDeployment(c *gin.Context) {
	ns := c.Query("namespace")
	deployment := c.Query("deployment")

	client := globalConfig.MyClient.Client
	err := client.ExtensionsV1beta1().Deployments(ns).Delete(deployment, &metav1.DeleteOptions{})
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
