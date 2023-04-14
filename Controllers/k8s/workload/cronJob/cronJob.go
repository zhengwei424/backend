package cronJob

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetCronJobsInfo(c *gin.Context) {
	var cronJobsInfo = make([]map[string]interface{}, 0)
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
	cronJobs, err := client.BatchV1beta1().CronJobs(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, cronJob := range cronJobs.Items {
		var cronJobInfo = make(map[string]interface{}, 0)
		cronJobInfo["name"] = cronJob.Name
		cronJobInfo["namespace"] = cronJob.Namespace
		cronJobInfo["labels"] = cronJob.Labels
		cronJobInfo["creationTimestamp"] = tools.DeltaTime(cronJob.CreationTimestamp.UTC(), time.Now())
		cronJobsInfo = append(cronJobsInfo, cronJobInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": cronJobsInfo,
		"msg":  "ok",
	})
}
