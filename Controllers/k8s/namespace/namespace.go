package namespace

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetNamespacesInfo(c *gin.Context) {
	namespacesInfo := make([]map[string]interface{}, 0)

	client := globalConfig.MyClient.Client
	opts := v1.ListOptions{}
	namespaces, err := client.CoreV1().Namespaces().List(opts)
	if err != nil {
		fmt.Println(err)
	}
	for _, ns := range namespaces.Items {
		var nsInfo = make(map[string]interface{}, 0)
		nsInfo["name"] = ns.Name
		nsInfo["status"] = ns.Status.Phase
		nsInfo["creationTimestamp"] = tools.DeltaTime(ns.CreationTimestamp.UTC(), time.Now())
		namespacesInfo = append(namespacesInfo, nsInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": namespacesInfo,
		"msg":  "ok",
	})
}
