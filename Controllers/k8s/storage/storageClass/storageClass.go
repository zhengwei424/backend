package storageClass

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetStorageClassesInfo(c *gin.Context) {
	var storageClassesInfo = make([]map[string]interface{}, 0)
	// 获取k8s自身资源管理客户端实例
	client := globalConfig.MyClient.Client

	opts := v1.ListOptions{}
	storageClasses, err := client.StorageV1().StorageClasses().List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, storageClass := range storageClasses.Items {
		var storageClassInfo = make(map[string]interface{}, 0)
		storageClassInfo["name"] = storageClass.Name
		storageClassInfo["labels"] = storageClass.Labels
		storageClassInfo["provisioner"] = storageClass.Provisioner
		storageClassInfo["reclaimPolicy"] = storageClass.ReclaimPolicy
		storageClassInfo["mountOptions"] = storageClass.MountOptions
		storageClassInfo["volumeBindingMode"] = storageClass.VolumeBindingMode
		storageClassInfo["creationTimestamp"] = tools.DeltaTime(storageClass.CreationTimestamp.UTC(), time.Now())
		storageClassesInfo = append(storageClassesInfo, storageClassInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": storageClassesInfo,
		"msg":  "ok",
	})
}
