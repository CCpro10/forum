package main

import (
	"forum/controller"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "net/http"
)

func RegisterHandle(c *gin.Context)  {

}
func main() {
	// 连接数据库
	//err := dao.InitMySQL()
	//if err != nil {
	//	panic(err)
	//}
	//defer dao.Close()  // 程序退出关闭数据库连接
	r:=gin.Default()

	r.POST("/register", RegisterHandle )
	r.POST("/login",controller.LoginHandle)

	Usergroup:=r.Group("/user",controller.JWTAuthMiddleware())
{
	Usergroup.GET("/home",func(c * gin.Context){
    username:=c.MustGet("username").(string)
    c.JSON(200,gin.H{
    	"code":2000,
		"msg":"message",
		"date":gin.H{"username":username},
	})
	})

}
    r.Run(":9999")
}