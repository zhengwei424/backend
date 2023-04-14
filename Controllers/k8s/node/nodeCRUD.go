package node

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func CreateNode(c *gin.Context) {
	nodeInfo := new(v1.Node)
	if err := c.BindJSON(nodeInfo); err == nil {
		fmt.Println(nodeInfo)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().Nodes().Create(nodeInfo)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		//"data": nodeInfo,
		"msg": "ok",
	})
}

func GetNode(c *gin.Context) {
	node := c.Query("node")

	nodeInfo := new(v1.Node)
	client := globalConfig.MyClient.Client
	nodeInfo, err := client.CoreV1().Nodes().Get(node, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		nodeInfo.Kind = "Node"
		nodeInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nodeInfo,
			"msg":  "ok",
		})
	}
}

func UpdateNode(c *gin.Context) {
	nodeInfo := new(v1.Node)
	if err := c.BindJSON(nodeInfo); err == nil {
		fmt.Println(*nodeInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().Nodes().Update(nodeInfo)
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

func DeleteNode(c *gin.Context) {
	node := c.Query("node")

	client := globalConfig.MyClient.Client
	err := client.CoreV1().Nodes().Delete(node, &metav1.DeleteOptions{})
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
