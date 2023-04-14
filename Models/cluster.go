package Models

import (
	"gorm.io/gorm"
	"time"
)

type Cluster struct {
	ID         uint64 `gorm:"column:id;primaryKey" json:"id"`
	APIVersion string `gorm:"column:apiVersion" json:"apiVersion"`
	Kind       string `gorm:"column:kind" json:"kind"`
	// context
	ContextName string `gorm:"column:contextName" json:"context-name"`
	Cluster     string `gorm:"column:cluster" json:"cluster,omitempty"`
	AuthInfo    string `gorm:"column:user" json:"user,omitempty"`
	// cluster
	Server                   string `gorm:"column:server" json:"server,omitempty"`
	InsecureSkipTLSVerify    bool   `gorm:"column:insecureSkipTLSVerify" json:"insecure-skip-tls-verify,omitempty"`
	CertificateAuthorityData string `gorm:"column:certificateAuthorityData" json:"certificate-authority-data,omitempty"`
	// user
	ClientCertificateData string `gorm:"column:clientCertificateData" json:"client-certificate-data,omitempty"`
	ClientKeyData         string `gorm:"column:clientKeyData" json:"client-key-data,omitempty"`
	Token                 string `gorm:"column:token" json:"token,omitempty"`
	// time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Insert 插入数据
func (cluster Cluster) Insert(db *gorm.DB) error {
	// 将user实例插入到数据库中
	//result := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").Create(user)
	result := db.Create(cluster)

	// 新建一条记录，为指定的字段分配值
	//db.Select("Username", "Name", "Email").Create(user)

	// 新建一条记录，忽略指定字段
	//db.Omit("Username", "Name", "Email").Create(user)

	return result.Error
}

// BatchInsertCluster 批量插入
func BatchInsertCluster(clusters []Cluster, db *gorm.DB) {
	// 1. 少量数据
	//var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
	//db.Create(&users)

	// 2. 大量数据可以指定每次批量的大小
	//var users = []User{{Name: "jinzhu_1"}, ...., {Name: "jinzhu_10000"}}
	// batch size 100
	db.CreateInBatches(clusters, 100)
}

// GetAllClusters 获取所有的集群信息
func GetAllClusters(db *gorm.DB) (clusters []Cluster, err error) {
	result := db.Model(&Cluster{}).Find(&clusters)
	return clusters, result.Error
}

// SpecifiedCluster 返回指定字段给前端，用于展示大致的集群信息，供用户筛选、连接
type SpecifiedCluster struct {
	ID          uint64 `gorm:"column:id" json:"id"`
	ContextName string `gorm:"column:contextName" json:"context-name"`
	Server      string `gorm:"column:server" json:"server,omitempty"`
	Cluster     string `gorm:"column:cluster" json:"cluster,omitempty"`
	AuthInfo    string `gorm:"column:user" json:"user,omitempty"`
}

// GetSpecifiedField 获取指定字段用于前端展示
func GetSpecifiedField(db *gorm.DB) (specifiedCluster []SpecifiedCluster, err error) {
	// 会返回指定字段之外的其他字段
	//result := db.Model(&Cluster{}).Select("id", "contextName", "cluster", "user").Find(&clusters)
	// 可以定义根据需要返回的字段，定义一个新的结构体，使用scan将查询结果返回到该结构体中
	result := db.Model(&Cluster{}).Select("id", "contextName", "server", "cluster", "user").Scan(&specifiedCluster)
	return specifiedCluster, result.Error
}

// 通过id和context返回指定集群信息，用于连接集群
func QueryCluster(id uint64, context string, db *gorm.DB) (cluster Cluster, err error) {
	result := db.Model(&Cluster{}).Where("id = ? and contextName = ?", id, context).First(&cluster)
	return cluster, result.Error
}
