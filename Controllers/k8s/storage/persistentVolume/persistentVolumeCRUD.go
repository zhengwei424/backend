package persistentVolume

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetPersistentVolume(c *gin.Context) {
	persistentVolume := c.Query("persistentVolume")
	persistentVolumeInfo := new(v1.PersistentVolume)

	client := globalConfig.MyClient.Client
	persistentVolumeInfo, err := client.CoreV1().PersistentVolumes().Get(persistentVolume, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		persistentVolumeInfo.Kind = "PersistentVolume"
		persistentVolumeInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": persistentVolumeInfo,
			"msg":  "ok",
		})
	}
}

func UpdatePersistentVolume(c *gin.Context) {
	persistentVolumeInfo := new(v1.PersistentVolume)
	if err := c.BindJSON(persistentVolumeInfo); err == nil {
		fmt.Println(*persistentVolumeInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().PersistentVolumes().Update(persistentVolumeInfo)
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

func DeletePersistentVolume(c *gin.Context) {
	persistentVolume := c.Query("persistentVolume")

	client := globalConfig.MyClient.Client
	err := client.CoreV1().PersistentVolumes().Delete(persistentVolume, &metav1.DeleteOptions{})
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
