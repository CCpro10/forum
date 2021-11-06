package main

import (
	"forum/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//帖子
type Article struct {
	Id  int
	Context  string
	Author  string
	Replies   []Reply//所有回复

}

//帖子下面的回复
type Reply struct{
	Id int
	Author  string
	Context string
	//用来存帖子下面回复的回复
	Comments  []Comment
}

//帖子下面回复的回复
type Comment struct {
	Author  string
	Context string
}
type User struct {
	Id   int
	Name  string  //用户名
   	Pwd   string

}

func main() {
	// 连接数据库
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	defer dao.Close()  // 程序退出关闭数据库连接
	r:=gin.Default()
	usergroup:=r.Group("/user",AuthMiddleWare())

	r.POST("/auth", authHandler)

}