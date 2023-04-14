package namespace

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func CreateNamespace(c *gin.Context) {
	client := globalConfig.MyClient.Client
	namespaceInfo := new(v1.Namespace)
	if err := c.BindJSON(namespaceInfo); err == nil {
		fmt.Println(namespaceInfo)
	}
	_, err := client.CoreV1().Namespaces().Create(namespaceInfo)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		//"data": namespaceInfo,
		"msg": "ok",
	})
}

func GetNamespace(c *gin.Context) {
	namespace := c.Query("namespace")
	namespaceInfo := new(v1.Namespace)

	client := globalConfig.MyClient.Client
	namespaceInfo, err := client.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		namespaceInfo.Kind = "Namespace"
		namespaceInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": namespaceInfo,
			"msg":  "ok",
		})
	}
}

func UpdateNamespace(c *gin.Context) {
	namespaceInfo := new(v1.Namespace)
	if err := c.BindJSON(namespaceInfo); err == nil {
		fmt.Println(*namespaceInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().Namespaces().Update(namespaceInfo)
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

func DeleteNamespace(c *gin.Context) {
	namespace := c.Query("namespace")
	client := globalConfig.MyClient.Client

	err := client.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
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
