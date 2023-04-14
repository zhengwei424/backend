package networkPolicy

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetNetworkPoliciesInfo(c *gin.Context) {
	var networkPoliciesInfo = make([]map[string]interface{}, 0)
	var qry, ns string
	// 获取k8s自身资源管理客户端实例
	client := globalConfig.MyClient.Client

	// 获取url中的namespace请求参数
	qry = c.Query("namespace")
	if qry == "all" {
		ns = ""
	} else {
		ns = qry
	}

	opts := v1.ListOptions{}
	networkPolicies, err := client.ExtensionsV1beta1().NetworkPolicies(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, networkPolicy := range networkPolicies.Items {
		var networkPolicyInfo = make(map[string]interface{}, 0)
		networkPolicyInfo["name"] = networkPolicy.Name
		networkPolicyInfo["namespace"] = networkPolicy.Namespace
		networkPolicyInfo["labels"] = networkPolicy.Labels
		networkPolicyInfo["creationTimestamp"] = tools.DeltaTime(networkPolicy.CreationTimestamp.UTC(), time.Now())
		networkPoliciesInfo = append(networkPoliciesInfo, networkPolicyInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": networkPoliciesInfo,
		"msg":  "ok",
	})
}
