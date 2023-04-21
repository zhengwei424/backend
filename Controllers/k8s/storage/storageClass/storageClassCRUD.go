package storageClass

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetStorageClass(c *gin.Context) {
	storageClass := c.Query("storageClass")
	storageClassInfo := new(v1.StorageClass)

	client := globalConfig.MyClient.Client
	storageClassInfo, err := client.StorageV1().StorageClasses().Get(storageClass, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		storageClassInfo.Kind = "StorageClass"
		storageClassInfo.APIVersion = "v1.k8s.io/v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": storageClassInfo,
			"msg":  "ok",
		})
	}
}

func UpdateStorageClass(c *gin.Context) {
	storageClassInfo := new(v1.StorageClass)
	if err := c.BindJSON(storageClassInfo); err == nil {
		fmt.Println(*storageClassInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.StorageV1().StorageClasses().Update(storageClassInfo)
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

func DeleteStorageClass(c *gin.Context) {
	storageClass := c.Query("storageClass")

	client := globalConfig.MyClient.Client
	err := client.StorageV1().StorageClasses().Delete(storageClass, &metav1.DeleteOptions{})
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
