package clusterManagement

import (
	"backend/Databases"
	"backend/Models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetClusterKubeConfigPartInfo(c *gin.Context) {
	var db = new(gorm.DB)
	var err error
	var clusters = make([]Models.SpecifiedCluster, 0)

	db, err = Databases.ConnMysql("backend")
	clusters, err = Models.GetSpecifiedField(db)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": clusters,
	})
}
