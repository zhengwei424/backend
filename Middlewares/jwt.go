package Middlewares

import (
	"backend/tools"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		token := c.GetHeader("token")
		// 判断请求头是否携带token字段
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "token must be passed",
			})
			//
			c.Abort()
			return
		}
		// 解析请求头的token字段，并做验证
		// tokenAfterParse中有很多信息可供有验证使用，以后研究？？？？？？？？？？、
		tokenAfterParse, claims, err := tools.ParseToken(token)
		if err != nil || !tokenAfterParse.Valid {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "parse token failed",
			})
			c.Abort()
			return
		}
		// 从token中解析出来的数据挂载到上下文上，方便后面的控制器使用
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
