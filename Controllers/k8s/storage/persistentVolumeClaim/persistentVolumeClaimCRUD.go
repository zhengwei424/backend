package persistentVolumeClaim

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetPersistentVolumeClaim(c *gin.Context) {
	ns := c.Query("namespace")
	persistentVolumeClaim := c.Query("persistentVolumeClaim")
	persistentVolumeClaimInfo := new(v1.PersistentVolumeClaim)

	client := globalConfig.MyClient.Client
	persistentVolumeClaimInfo, err := client.CoreV1().PersistentVolumeClaims(ns).Get(persistentVolumeClaim, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		persistentVolumeClaimInfo.Kind = "PersistentVolumeClaim"
		persistentVolumeClaimInfo.APIVersion = "v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": persistentVolumeClaimInfo,
			"msg":  "ok",
		})
	}
}

func UpdatePersistentVolumeClaim(c *gin.Context) {
	persistentVolumeClaimInfo := new(v1.PersistentVolumeClaim)
	if err := c.BindJSON(persistentVolumeClaimInfo); err == nil {
		fmt.Println(*persistentVolumeClaimInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.CoreV1().PersistentVolumeClaims(persistentVolumeClaimInfo.Namespace).Update(persistentVolumeClaimInfo)
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

func DeletePersistentVolumeClaim(c *gin.Context) {
	ns := c.Query("namespace")
	persistentVolumeClaim := c.Query("persistentVolumeClaim")

	client := globalConfig.MyClient.Client
	err := client.CoreV1().PersistentVolumeClaims(ns).Delete(persistentVolumeClaim, &metav1.DeleteOptions{})
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
