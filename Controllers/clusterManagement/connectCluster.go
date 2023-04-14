package clusterManagement

import (
	"backend/Databases"
	"backend/Models"
	"backend/globalConfig"
	"backend/kubeconfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type contextInfo struct {
	ID          uint64 `json:"id"`
	ContextName string `json:"context-name"`
}

func ConnectCluster(c *gin.Context) {
	var ctx = new(contextInfo)
	var err error

	// 将post请求参数转换为ctx
	err = c.BindJSON(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// 连接到backend数据库
	var db = new(gorm.DB)
	db, err = Databases.ConnMysql("backend")
	if err != nil {
		panic(err)
	}

	// 根据ctx信息，查询数据库，得到完整的Models.Cluster结构体信息
	var cluster Models.Cluster
	cluster, err = Models.QueryCluster(ctx.ID, ctx.ContextName, db)
	if err != nil {
		panic(err)
	}

	// 解析查询结果Models.Cluster为KubeConfig类型
	var kc kubeconfig.KubeConfig
	kc = kubeconfig.ParseClusterToKubeConfig(cluster)

	// 通过KubeConfig生成rest.Config
	globalConfig.MyCfg, err = kc.ParseKubeConfigToRestConfig(ctx.ContextName)
	if err != nil {
		panic(err)
	}

	// 生成Clientset
	globalConfig.MyClient, err = tools.GenerateClient(globalConfig.MyCfg)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}
