package Middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

// recover 错误,转string
func errorToString(err interface{}) string {
	switch v := err.(type) {
	case error:
		return v.Error()
	default:
		return err.(string)
	}
}

// Recovery 全局统一异常处理
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印错误堆栈信息
				debug.PrintStack()
				// 返回Json
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  errorToString(err),
				})
				// 终止后续接口的调用
				c.Abort()
			}
		}()
		// 加载完defer recover,继续后续接口调用
		c.Next()
	}
}
