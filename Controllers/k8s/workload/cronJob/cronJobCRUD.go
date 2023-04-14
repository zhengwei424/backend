package cronJob

import (
	"backend/globalConfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetCronJob(c *gin.Context) {
	ns := c.Query("namespace")
	cronJob := c.Query("cronJob")
	cronJobInfo := new(v1beta1.CronJob)

	client := globalConfig.MyClient.Client
	cronJobInfo, err := client.BatchV1beta1().CronJobs(ns).Get(cronJob, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": nil,
			"msg":  err,
		})
	} else {
		cronJobInfo.Kind = "CronJob"
		cronJobInfo.APIVersion = "batch/v1beta1"
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": cronJobInfo,
			"msg":  "ok",
		})
	}
}

func UpdateCronJob(c *gin.Context) {
	cronJobInfo := new(v1beta1.CronJob)
	if err := c.BindJSON(cronJobInfo); err == nil {
		fmt.Println(*cronJobInfo)
	} else {
		fmt.Println(err)
	}

	client := globalConfig.MyClient.Client
	_, err := client.BatchV1beta1().CronJobs(cronJobInfo.Namespace).Update(cronJobInfo)
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
	cronJob := c.Query("cronJob")

	client := globalConfig.MyClient.Client
	err := client.BatchV1beta1().CronJobs(ns).Delete(cronJob, &metav1.DeleteOptions{})
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
