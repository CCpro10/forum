package main

import (

	"forum/models"
	"github.com/gin-gonic/gin"
	"forum/dao"
	"net/http"
)

func main() {
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	r:=gin.Default()
	r.POST("/aaa", func(c *gin.Context) {
		var forum models.Forum
		c.ShouldBind(&forum)
		dao.DB.AutoMigrate(models.Forum{})
		dao.DB.Create(&forum)
		c.JSON(http.StatusOK,gin.H{
			"msg":"ok",
			"data": forum,
		})
	})
    r.Run(":9999")
}
