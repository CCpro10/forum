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
	//展示某个论坛的posts
	Usergroup.GET("/posts",controller.Showposts)
	//展示某个post的comments
	Usergroup.GET("/comments",controller.Showcomments)
	//展示某个comment的replies
	Usergroup.GET("/replies",controller.Showreplies)

	//home是用户主页会展示个人信息
	Usergroup.GET("/home",controller.ShowHomePage)
	//
	Usergroup.PUT("/username",controller.ChangeUsername)

	//发帖
	Usergroup.POST("/post",controller.CreatPost)
	//发帖子评论
	Usergroup.POST("/comment",controller.CreatComment)
	//回复帖子的评论
	Usergroup.POST("/reply",controller.CreatReply)

	//设置论坛发帖权限,传入forumcod和postpermission
	Usergroup.PUT("/postpermission",controller.SetPostPermission)
	//设置论坛发帖权限,传入forumcod和commentpermission
	Usergroup.PUT("/commentpermission",controller.SetCommentPermission)
	//设置论坛发帖权限,传入forumcod和accesspermission
	Usergroup.PUT("/accesspermission",controller.SetAccessPermission)
}
	//写一个开发者的路由分组
	Developergroup:=r.Group("/developer",controller.JWTAuthMiddleware())
{
	//设置管理员
	Developergroup.POST("/managerlist",controller.SetManager)
}



    r.Run(":9999")
}