package main

import (

	"forum/models"
	"github.com/gin-gonic/gin"
	"forum/dao"
	"net/http"
)

//初始化forumpermission
func main() {
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	r:=gin.Default()
	r.POST("/init", func(c *gin.Context) {
		var forumpermission models.Forumpermission
		_= c.ShouldBind(&forumpermission)
		dao.DB.AutoMigrate(models.Forumpermission{})
		dao.DB.Create(&forumpermission)
		c.JSON(http.StatusOK,gin.H{
			"msg":  "ok",
			"data": forumpermission,
		})
	})
    r.Run(":9999")
}
