package pod

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
	"time"
)

func GetPodsInfo(c *gin.Context) {
	var podsInfo = make([]map[string]interface{}, 0)
	var qry, ns string
	// 获取k8s自身资源管理客户端实例
	client := globalConfig.MyClient.Client

	// 获取uri中传入的namespace
	qry = c.Query("namespace")
	if qry == "all" {
		ns = ""
	} else {
		ns = qry
	}

	var opts v1.ListOptions
	pods, err := client.CoreV1().Pods(ns).List(opts)
	if err != nil {
		panic(err)
	}

	for _, pod := range pods.Items {
		// 获取单个pod的信息
		var podInfo = make(map[string]interface{}, 0)
		podInfo["name"] = pod.Name
		podInfo["namespace"] = pod.Namespace
		podInfo["annotations"] = pod.Annotations
		podInfo["ower"] = pod.ObjectMeta.OwnerReferences
		if len(pod.ObjectMeta.OwnerReferences) != 0 {
			podInfo["controlledBy"] = pod.ObjectMeta.OwnerReferences[0].Kind
		} else {
			podInfo["controlledBy"] = nil
		}
		podInfo["uid"] = pod.UID
		podInfo["labels"] = pod.Labels
		podInfo["nodeIP"] = pod.Status.HostIP
		podInfo["status"] = pod.Status.Phase
		podInfo["qos"] = pod.Status.QOSClass
		podInfo["priorityClass"] = pod.Spec.PriorityClassName
		podInfo["priority"] = pod.Spec.Priority
		podInfo["conditions"] = pod.Status.Conditions
		podInfo["tolerations"] = pod.Spec.Tolerations
		podInfo["podIP"] = pod.Status.PodIP
		podInfo["restartCount"] = int32(0)
		podInfo["startTime"] = tools.DeltaTime(pod.Status.StartTime.UTC(), time.Now())
		podInfo["created"] = fmt.Sprintf("%s(%s)", tools.DeltaTime(pod.Status.StartTime.UTC(), time.Now()), pod.Status.StartTime.UTC().String())
		// 获取某个pod中所有container的信息
		var containersInfo = make([]map[string]interface{}, 0)
		for _, container := range pod.Spec.Containers {
			//portLength := len(container.Ports)
			//if portLength > 1 {
			//	fmt.Println("+++++++++++++++++++++++")
			//	fmt.Println(pod.Namespace)
			//	fmt.Println(pod.Name)
			//	fmt.Println(container.Ports)
			//	fmt.Println("+++++++++++++++++++++++")
			//}
			// 获取某个pod中单个container的信息
			var containerInfo = make(map[string]interface{}, 0)
			containerInfo["name"] = container.Name
			containerInfo["image"] = container.Image
			containerInfo["requests"] = map[string]resource.Quantity{
				"cpu":    container.Resources.Requests["cpu"],
				"memory": container.Resources.Requests["memory"],
			}
			containerInfo["limits"] = map[string]resource.Quantity{
				"cpu":    container.Resources.Limits["cpu"],
				"memory": container.Resources.Limits["memory"],
			}
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.RestartCount > podInfo["restartCount"].(int32) {
					podInfo["restartCount"] = cs.RestartCount
				}
				if containerInfo["name"] == cs.Name {
					// 重启次数
					containerInfo["restartCount"] = cs.RestartCount
					if cs.ContainerID != "" {
						containerInfo["containerID"] = strings.Split(cs.ContainerID, "//")[1][0:12]
					}
					// containerState
					// 结构体判断
					if !tools.IsStructEmpty(cs.State, v12.ContainerState{}) {
						// 结构体指针判断
						if cs.State.Running != nil {
							containerInfo["state"] = "Running"
							containerInfo["startedAt"] = tools.UTCTimeToTimeStr(cs.State.Running.StartedAt.UTC(), time.RFC3339)
						} else if cs.State.Waiting != nil {
							containerInfo["state"] = "Waiting"
							containerInfo["reason"] = cs.State.Waiting.Reason
							containerInfo["message"] = cs.State.Waiting.Message
						} else if cs.State.Terminated != nil {
							containerInfo["state"] = "Terminated"
							containerInfo["reason"] = cs.State.Terminated.Reason
							containerInfo["message"] = cs.State.Terminated.Message
						}
					}
					// containerLastState
					containerInfo["lastTerminationState"] = map[string]string{}
					if cs.LastTerminationState.Terminated != nil {
						containerInfo["lastTerminationState"] = map[string]string{
							"reason":      cs.LastTerminationState.Terminated.Reason,
							"message":     cs.LastTerminationState.Terminated.Message,
							"containerID": strings.Split(cs.LastTerminationState.Terminated.ContainerID, "//")[1][0:12],
							"startedAt":   cs.LastTerminationState.Terminated.StartedAt.UTC().String(),
							"finishedAt":  cs.LastTerminationState.Terminated.FinishedAt.UTC().String(),
						}
					}
				}
			}
			containersInfo = append(containersInfo, containerInfo)
		}
		podInfo["containers"] = containersInfo
		podsInfo = append(podsInfo, podInfo)
	}
	//tmp, _ := json.MarshalIndent(podsInfo, "", " ")
	//fmt.Println(string(tmp))
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": podsInfo,
		"msg":  "ok",
	})
}
