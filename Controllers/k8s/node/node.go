package node

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetNodesInfo(c *gin.Context) {
	nodesInfo := make([]map[string]interface{}, 0)
	// 获取k8s自身资源管理客户端实例
	client := globalConfig.MyClient.Client

	opts := v1.ListOptions{}
	nodes, err := client.CoreV1().Nodes().List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, node := range nodes.Items {
		var nodeInfo = make(map[string]interface{}, 0)
		nodeInfo["name"] = node.Name
		nodeInfo["labels"] = node.Labels
		nodeInfo["Annotations"] = node.Annotations
		nodeInfo["creationTimestamp"] = tools.DeltaTime(node.CreationTimestamp.UTC(), time.Now())
		nodeInfo["taints"] = node.Spec.Taints
		nodeInfo["address"] = map[string]string{"InternalIP": "", "Hostname": ""}
		for _, v := range node.Status.Addresses {
			fmt.Printf("%v", v)
			if v.Type == "InternalIP" {
				// interface需要断言自己需要的类型
				nodeInfo["address"].(map[string]string)["InternalIP"] = v.Address
			}
			if v.Type == "Hostname" {
				nodeInfo["address"].(map[string]string)["Hostname"] = v.Address
			}
		}
		nodeInfo["os"] = fmt.Sprintf("%s(%s)", node.Status.NodeInfo.OperatingSystem, node.Status.NodeInfo.Architecture)
		nodeInfo["osImage"] = node.Status.NodeInfo.OSImage
		nodeInfo["kernelVersion"] = node.Status.NodeInfo.KernelVersion
		nodeInfo["kubeletVersion"] = node.Status.NodeInfo.KubeletVersion
		nodeInfo["containerRuntimeVersion"] = node.Status.NodeInfo.ContainerRuntimeVersion
		nodeInfo["allocatable"] = node.Status.Allocatable
		nodeInfo["capacity"] = node.Status.Capacity
		//nodeInfo["allocatable_cpu"] = node.Status.Allocatable.Cpu().Value()
		//nodeInfo["allocatable_mem"] = node.Status.Allocatable.Memory().Value()
		//nodeInfo["capacity_cpu"] = node.Status.Capacity.Cpu().Value()
		//nodeInfo["capacity_mem"] = node.Status.Capacity.Memory().Value()
		nodeInfo["conditions"] = node.Status.Conditions
		nodesInfo = append(nodesInfo, nodeInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nodesInfo,
		"msg":  "ok",
	})
}
