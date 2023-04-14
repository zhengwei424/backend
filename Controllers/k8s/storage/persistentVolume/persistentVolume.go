package persistentVolume

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetPersistentVolumesInfo(c *gin.Context) {
	var persistentVolumesInfo = make([]map[string]interface{}, 0)
	// 获取k8s自身资源管理客户端实例
	client := globalConfig.MyClient.Client

	opts := v1.ListOptions{}
	persistentVolumes, err := client.CoreV1().PersistentVolumes().List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, persistentVolume := range persistentVolumes.Items {
		var persistentVolumeInfo = make(map[string]interface{}, 0)
		persistentVolumeInfo["name"] = persistentVolume.Name
		persistentVolumeInfo["labels"] = persistentVolume.Labels
		persistentVolumeInfo["storage"] = persistentVolume.Spec.Capacity.StorageEphemeral().Value() // 容量
		persistentVolumeInfo["accessModes"] = persistentVolume.Spec.AccessModes
		persistentVolumeInfo["persistentVolumeReclaimPolicy"] = persistentVolume.Spec.PersistentVolumeReclaimPolicy
		persistentVolumeInfo["mountOptions"] = persistentVolume.Spec.MountOptions
		persistentVolumeInfo["volumeMode"] = persistentVolume.Spec.VolumeMode
		persistentVolumeInfo["creationTimestamp"] = tools.DeltaTime(persistentVolume.CreationTimestamp.UTC(), time.Now())
		persistentVolumesInfo = append(persistentVolumesInfo, persistentVolumeInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": persistentVolumesInfo,
		"msg":  "ok",
	})
}
