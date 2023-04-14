package Middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 解决跨域问题
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 在响应头中展示：在请求头中可以设置哪些字段（如：在请求头中加入token字段）
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding,X-Token, token, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 解决axios无法接收到c.DataFromReader的extraHeader问题
		// 设置之后，才能在axios中接收到filename字段，否则即使通过F12调试能看到响应头中的filenamne，但是在axios的then中也无法接收到filename字段
		c.Writer.Header().Set("Access-Control-Expose-Headers", "filename")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, UPDATE")

		// 处理options嗅探请求
		if c.Request.Method == "OPTIONS" {
			// 可以拦截，也可以response，无所谓
			c.AbortWithStatus(http.StatusNoContent)
			//c.JSON(http.StatusOK, "Options Request!")
		}
		// 继续处理请求
		c.Next()
	}
}
