package Models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint64 `gorm:"primaryKey"`
	Username    string `gorm:"column:username;size:20;not null;unique" json:"username"`
	Password    string `gorm:"column:password;size:256;not null" json:"password"`
	Name        string `gorm:"column:name;not null" json:"name"`
	Avatar      string `gorm:"column:avatar;" json:"avatar"`
	PhoneNumber string `gorm:"column:phone_number;size:11;not null;unique" json:"phone_number"`
	Email       string `gorm:"column:email;size:50;not null;unique" json:"email"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Insert 插入数据
func (user User) Insert(db *gorm.DB) error {
	// 将user实例插入到数据库中
	//result := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").Create(user)
	result := db.Create(user)

	// 新建一条记录，为指定的字段分配值
	//db.Select("Username", "Name", "Email").Create(user)

	// 新建一条记录，忽略指定字段
	//db.Omit("Username", "Name", "Email").Create(user)

	return result.Error
}

// BatchInsertUser 批量插入
func BatchInsertUser(users []User, db *gorm.DB) {
	// 1. 少量数据
	//var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
	//db.Create(&users)

	// 2. 大量数据可以指定每次批量的大小
	//users = []User{{Name: "jinzhu_1"}, ...., {Name: "jinzhu_10000"}}
	// batch size 100
	db.CreateInBatches(users, 100)
}

// IsExist 判断用户名是否存在
func IsExist(username string, db *gorm.DB) (user User, err error) {
	// 执行完之后，查询结果就在user里
	result := db.Model(&User{}).Where("username=?", username).Find(&user)

	return user, result.Error
}
