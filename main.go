package main

import (
	"forum/controller"
	"forum/dao"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	//连接数据库
	err := dao.InitMySQL()
	//dao.DB.AutoMigrate(models.UserInfo{})
	if err != nil {
		panic(err)
	}
	//defer dao.Close()  // 程序退出关闭数据库连接

	r:=gin.Default()
	r.POST("/register", controller.RegisterHandle )
	r.POST("/login",controller.LoginHandle)




	Usergroup:=r.Group("/user",controller.JWTAuthMiddleware())
{
	Usergroup.GET("/home",func(c * gin.Context){
    userid:=c.MustGet("userid")
    c.JSON(200,gin.H{
    	"code":2000,
		"msg":"访问成功",
		"date":gin.H{"userid":userid},
	})
	})
	Usergroup.POST("/post",controller.Postcreat)


}
    r.Run(":9999")
}