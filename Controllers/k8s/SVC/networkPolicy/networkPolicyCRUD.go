package networkPolicy

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetNetworkPolicy(c *gin.Context) {
	ns := c.Query("namespace")
	networkPolicy := c.Query("networkPolicy")
	networkPolicyInfo := new(v1beta1.NetworkPolicy)

	client := globalConfig.MyClient.Client
	networkPolicyInfo, err := client.ExtensionsV1beta1().NetworkPolicies(ns).Get(networkPolicy, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		networkPolicyInfo.Kind = "NetworkPolicy"
		networkPolicyInfo.APIVersion = "extensions/v1beta1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": networkPolicyInfo,
			"msg":  "ok",
		})
	}
}

func UpdateNetworkPolicy(c *gin.Context) {
	networkPolicyInfo := new(v1beta1.NetworkPolicy)
	if err := c.BindJSON(networkPolicyInfo); err == nil {
		fmt.Println(*networkPolicyInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.ExtensionsV1beta1().NetworkPolicies(networkPolicyInfo.Namespace).Update(networkPolicyInfo)
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

func DeleteNetworkPolicy(c *gin.Context) {
	ns := c.Query("namespace")
	networkPolicy := c.Query("networkPolicy")

	client := globalConfig.MyClient.Client
	err := client.ExtensionsV1beta1().NetworkPolicies(ns).Delete(networkPolicy, &metav1.DeleteOptions{})
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
