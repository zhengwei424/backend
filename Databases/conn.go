package Databases

import (
	"backend/tools"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func ConnMysql(database string) (*gorm.DB, error) {
	var err error
	var db *sql.DB
	config := tools.GetConfig().MysqlConfig
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		database,
	)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("connError:", err)
		panic(err)
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println("connError1:", err)
		panic(err)
	}

	// 非连接池，需要手动关闭数据库
	//// 获取通用数据库对象 sql.DB，然后使用其提供的功能
	//sqlDB, err := gormDB.DB()
	//
	//// Ping
	//sqlDB.Ping()
	//
	//// Close
	//sqlDB.Close()
	//
	//// 返回数据库统计信息
	//sqlDB.Stats()

	// 初始化连接池
	// 获取通用数据库对象，然后使用其提供的功能
	sqlDB, err := gormDB.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	return gormDB, err
}
