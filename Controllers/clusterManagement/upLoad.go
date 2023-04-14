package clusterManagement

import (
	"backend/Databases"
	"backend/Models"
	"backend/kubeconfig"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"net/http"
	"sigs.k8s.io/yaml"
)

func UpLoadKubeConfigFile(c *gin.Context) {
	var err error
	var fh = new(multipart.FileHeader)
	var db = new(gorm.DB)
	var tmpFile multipart.File

	fh, err = c.FormFile("file")
	if err != nil {
		log.Panicf("读取上传的kubeconfig文件失败: %s", err.Error())
	}

	tmpFile, err = fh.Open()
	if err != nil {
		log.Panic(err)
	}

	// 数组需要用make初始化长度，否则获取不到数据
	var content = make([]byte, fh.Size)

	_, err = tmpFile.Read(content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	defer func() {
		err = tmpFile.Close()
	}()
	// 利用反射读取multipart.FileHeader内的私有属性
	//file := reflect.ValueOf(*fh)
	// FieldByName的调用者必须时结构体，而不是指针
	//b := file.FieldByName("content").Bytes()
	//fmt.Println(string(b))

	// yaml []byte 转json []byte
	var kcfg []byte
	kcfg, err = yaml.YAMLToJSON(content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	var clusters = make([]Models.Cluster, 0)
	var kc = new(kubeconfig.KubeConfig)

	kc, err = kubeconfig.ParseByteToKubeConfig(kcfg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	clusters = kc.ParseKubeConfigToCluster()

	// 连接数据库
	db, err = Databases.ConnMysql("backend")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	// 插入clusters到数据库
	Models.BatchInsertCluster(clusters, db)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}
