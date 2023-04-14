package resourceTemplate

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func GetResourceTemplate(c *gin.Context) {
	// 初始化资源类型
	var rs string
	// 获取url中的resource请求参数
	rs = c.Query("resource")
	// 获取项目根路径
	rsPath := "./Controllers/k8s/other/resourceTemplate/templateFiles/" + rs + ".yaml"
	buf, err := ioutil.ReadFile(rsPath)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": string(buf),
	})
}
