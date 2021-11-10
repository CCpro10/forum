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
	r.POST("/register", controller.Register)
	r.POST("/login",controller.Login)
	r.POST("/developerlogin",controller.Developerlogin)


	Usergroup:=r.Group("/user",controller.JWTAuthMiddleware())
{
	//home是用户主页会展示个人信息
	Usergroup.GET("/home",controller.ShowHomePage)
	//发帖
	Usergroup.POST("/post",controller.CreatPost)
	//发帖子评论
	Usergroup.POST("/comment",controller.CreatComment)
	//发帖子评论的回复
	Usergroup.POST("/reply",controller.CreatReply)
}
	//写一个开发者的路由分组
	Developergroup:=r.Group("/manager",controller.JWTAuthMiddleware())
{
	Developergroup.POST("/managerlist",controller.SetManager)



}



    r.Run(":9999")
}