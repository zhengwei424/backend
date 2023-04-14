package Databases

import (
	"backend/Models"
)

// package中的init方法，默认会在import时被调用
func init() {
	// ---------------------------
	// 存放各类数据库的连接
	// --------------------------
	// 初始化连接
	backendUserConn, err := ConnMysql("backend")
	if err != nil {
		panic(err)
	}
	// 在backend库中初始化Models中的结构体（表）
	err = backendUserConn.AutoMigrate(&Models.User{}, &Models.Cluster{})
	if err != nil {
		panic(err)
	}
}
