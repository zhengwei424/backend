package job

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetJob(c *gin.Context) {
	ns := c.Query("namespace")
	job := c.Query("job")
	jobInfo := new(v1.Job)

	client := globalConfig.MyClient.Client
	jobInfo, err := client.BatchV1().Jobs(ns).Get(job, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		jobInfo.Kind = "Job"
		jobInfo.APIVersion = "batch/v1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": jobInfo,
			"msg":  "ok",
		})
	}
}

func UpdateJob(c *gin.Context) {
	jobInfo := new(v1.Job)
	if err := c.BindJSON(jobInfo); err == nil {
		fmt.Println(*jobInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.BatchV1().Jobs(jobInfo.Namespace).Update(jobInfo)
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

func DeleteJob(c *gin.Context) {
	ns := c.Query("namespace")
	job := c.Query("job")

	client := globalConfig.MyClient.Client
	err := client.BatchV1().Jobs(ns).Delete(job, &metav1.DeleteOptions{})
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
