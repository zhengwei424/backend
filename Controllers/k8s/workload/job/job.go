package job

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetJobsInfo(c *gin.Context) {
	var jobsInfo = make([]map[string]interface{}, 0)
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
	jobs, err := client.BatchV1().Jobs(ns).List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, job := range jobs.Items {
		var jobInfo = make(map[string]interface{}, 0)
		jobInfo["name"] = job.Name
		jobInfo["namespace"] = job.Namespace
		jobInfo["labels"] = job.Labels
		jobInfo["creationTimestamp"] = tools.DeltaTime(job.CreationTimestamp.UTC(), time.Now())
		jobsInfo = append(jobsInfo, jobInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": jobsInfo,
		"msg":  "ok",
	})
}
