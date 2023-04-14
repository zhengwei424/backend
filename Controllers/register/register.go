package register

import (
	"backend/Databases"
	"backend/Models"
	"backend/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type RegInfo struct {
	Username    string `gorm:"column:username;size:20;not null;unique" json:"username"`
	Password    string `gorm:"column:password;size:256;not null" json:"password"`
	Name        string `gorm:"column:name;size:20;not null" json:"name"`
	PhoneNumber string `gorm:"column:phone_number;size:11;not null;unique" json:"phone_number"`
	Email       string `gorm:"column:email;size:50;not null;unique" json:"email"`
}

func Register(c *gin.Context) {
	// 用户注册信息
	var regInfo = new(RegInfo)
	var err error
	var user Models.User
	var db = new(gorm.DB)
	var hashPassword string

	// 接收客户端发来的注册信息
	err = c.BindJSON(regInfo)
	if err != nil {
		log.Panicf("注册信息解析错误: %s", err.Error())
	}

	// 判断用户是否存在
	db, err = Databases.ConnMysql("backend")
	if err != nil {
		log.Panicf("MySQL连接backend库错误: %s", err)
	}
	user, err = Models.IsExist(regInfo.Username, db)
	if err != nil {
		log.Panicf("用户名查询错误: %s", err)
	}
	if (user != Models.User{}) {
		c.JSON(http.StatusOK, gin.H{
			"code": tools.UserExisted,
			"msg":  "username already exists",
		})
		return
	}
	// 用户不存在,可以注册，并返回加密密码
	hashPassword, err = tools.CryptPassword(regInfo.Password)
	if err != nil {
		log.Panicf("用户密码加密失败: %s", err.Error())
	}
	// 赋值给user实例，并调用insert方法插入数据库
	var regUser = Models.User{
		Username:    regInfo.Username,
		Password:    hashPassword,
		Name:        regInfo.Name,
		PhoneNumber: regInfo.PhoneNumber,
		Email:       regInfo.Email,
	}
	err = regUser.Insert(db)
	if err != nil {
		log.Panicf("MySQL插入注册信息错误: %s", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "register success",
	})
}
