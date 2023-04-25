package endpoint

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetEndpointsInfo(c *gin.Context) {
	var endpointsInfo = make([]map[string]interface{}, 0)
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
	endpoints, err := client.CoreV1().Endpoints(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, endpoint := range endpoints.Items {
		var endpointInfo = make(map[string]interface{}, 0)
		endpointInfo["name"] = endpoint.Name
		endpointInfo["namespace"] = endpoint.Namespace
		endpointInfo["labels"] = endpoint.Labels
		eps := make([]string, 0)
		for _, item1 := range endpoint.Subsets {
			epPorts := make([]string, 0)
			ep := make([]string, 0)
			if len(item1.Ports) > 0 {
				for _, n := range item1.Ports {
					epPorts = append(epPorts, fmt.Sprintf("%d/%s", n.Port, n.Protocol))
				}
			}
			if len(item1.Addresses) > 0 {
				if len(epPorts) > 0 {
					for _, i := range item1.Addresses {
						for _, j := range epPorts {
							ep = append(ep, fmt.Sprintf("%s:%s", i.IP, j))
						}
					}
				}
			}
			eps = append(eps, ep...)
		}
		endpointInfo["endpoints"] = eps
		endpointInfo["creationTimestamp"] = tools.DeltaTime(endpoint.CreationTimestamp.UTC(), time.Now())
		endpointsInfo = append(endpointsInfo, endpointInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": endpointsInfo,
		"msg":  "ok",
	})
}
