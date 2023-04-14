package persistentVolumeClaim

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetPersistentVolumeClaimsInfo(c *gin.Context) {
	var persistentVolumeClaimsInfo = make([]map[string]interface{}, 0)
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
	persistentVolumeClaims, err := client.CoreV1().PersistentVolumeClaims(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, persistentVolumeClaim := range persistentVolumeClaims.Items {
		var persistentVolumeClaimInfo = make(map[string]interface{}, 0)
		persistentVolumeClaimInfo["name"] = persistentVolumeClaim.Name
		persistentVolumeClaimInfo["labels"] = persistentVolumeClaim.Labels
		persistentVolumeClaimInfo["storage"] = persistentVolumeClaim.Spec.Resources.Requests.StorageEphemeral().Value() // 容量
		persistentVolumeClaimInfo["accessModes"] = persistentVolumeClaim.Spec.AccessModes
		persistentVolumeClaimInfo["volumeName"] = persistentVolumeClaim.Spec.VolumeName
		persistentVolumeClaimInfo["storageClassName"] = persistentVolumeClaim.Spec.StorageClassName
		persistentVolumeClaimInfo["volumeMode"] = persistentVolumeClaim.Spec.VolumeMode
		persistentVolumeClaimInfo["creationTimestamp"] = tools.DeltaTime(persistentVolumeClaim.CreationTimestamp.UTC(), time.Now())
		persistentVolumeClaimsInfo = append(persistentVolumeClaimsInfo, persistentVolumeClaimInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": persistentVolumeClaimsInfo,
		"msg":  "ok",
	})
}
