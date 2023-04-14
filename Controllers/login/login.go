package login

import (
	"backend/Databases"
	"backend/Models"
	"backend/tools"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type VerifyInfo struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	CaptchaID     string `json:"id"`
	CaptchaAnswer string `json:"answer"`
}

func Login(c *gin.Context) {
	// 获取axios的post请求参数
	var verifyInfo = new(VerifyInfo)
	var db = new(gorm.DB)
	var err error
	var user Models.User
	var token string

	// 解析登录参数
	err = c.BindJSON(verifyInfo)
	if err != nil {
		log.Panicf("登录参数解析错误: %s", err.Error())
	}

	//// 设置cookies，必须启用https
	//c.SetSameSite(http.SameSiteNoneMode)
	//c.SetCookie("vue_admin_template_token", "loginSuccess", 3600, "/", c.GetHeader("Origin"), true, false)

	// 先验证验证码，在验证用户名和密码
	// 1. 判断验证码是否正确
	ca := tools.AnswerCaptcha{
		ID:     verifyInfo.CaptchaID,
		Answer: verifyInfo.CaptchaAnswer,
		Clear:  false,
	}
	ok := ca.CaptchaVerify()
	if !ok {
		//c.JSON(http.StatusOK, gin.H{
		//	"code": tools.CaptchaError,
		//	"msg":  "captcha verify failed",
		//})
		c.JSON(http.StatusOK, gin.H{
			"code": tools.CaptchaError,
			"msg":  "captcha verify failed",
		})
		return
	}
	// 2. 判断用户名是否和数据库匹配（数据库中是否存在该用户）
	db, err = Databases.ConnMysql("backend")
	if err != nil {
		log.Panicf("MySQL连接backend库错误: %s", err)
	}
	user, err = Models.IsExist(verifyInfo.Username, db)
	if (user == Models.User{} || err != nil) {
		c.JSON(http.StatusOK, gin.H{
			"code": tools.UserOrPasswordError,
			"msg":  "username is not exist",
		})
		return
	}
	// 3. 判断密码是否与库中对应的密码相同
	result := tools.DeCryptPassword(user.Password, verifyInfo.Password)
	if !result {
		c.JSON(http.StatusOK, gin.H{
			"code": tools.UserOrPasswordError,
			"msg":  "incorrect password",
		})
		return
	}
	// 4. 生成jwt-token
	// 定义过期时间，7天后过期
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claim := tools.MyClaims{
		UserID:   user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),     // 发布时间
			Subject:   "token",               // 主题
			Issuer:    user.Username,         // 发布者
		},
	}

	token, err = claim.GenerateToken()
	if err != nil {
		log.Panicf("jwt-token生成失败: %s", err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "login success",
		"token": token,
	})

}
