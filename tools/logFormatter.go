package tools

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// MyLogFormatter 请求日志输出格式
func MyLogFormatter(params gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] %s %s %s %d %s \"%s\" %s\n",
		params.ClientIP,                       //请求客户端的IP地址
		params.TimeStamp.Format(time.RFC1123), //请求时间
		params.Method,                         //请求方法
		params.Path,                           //路由路径
		params.Request.Proto,                  //请求协议
		params.StatusCode,                     //http响应码
		params.Latency,                        //请求到响应的延时
		params.Request.UserAgent(),            //客户端代理程序
		params.ErrorMessage,                   //如果有错误,也打印错误信息
	)
}
