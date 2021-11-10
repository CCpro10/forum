package controller

import (
	_ "fmt"
	"forum/dao"
	"forum/models"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm"
	_ "golang.org/x/crypto/bcrypt"
	_ "log"
	"net/http"
)

//发布post
func CreatPost(c *gin.Context ) {
	var requestPost models.Post
	// 数据验证
	if err := c.ShouldBind(&requestPost); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "用户参数绑定失败:" + err.Error()})
	}
	if requestPost.Context==""||requestPost.Article==""{
		c.JSON(http.StatusUnprocessableEntity,gin.H{"msg":"帖子必须含标题和内容,帖子发布失败"})
		return
	}
	// 获取登录用户的id
	userid, _ := c.Get("userid")
	requestPost.Userid=userid.(int)

	//获取登录用户的名字作为post的作者
	var user models.UserInfo
	dao.DB.Where("ID=?",userid.(int)).First(&user)
	requestPost.Author=user.Username

	// 创建数据
	dao.DB.AutoMigrate(models.Post{})
	if err := dao.DB.Create(&requestPost).Error; err != nil {
		panic(err)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"mes":"创建成功",
		"data":requestPost,
		})
	return
}


//发布comment
func CreatComment(c *gin.Context ) {
	var requestComment models.Comment
	// 数据验证
	if err := c.ShouldBind(&requestComment); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "用户参数绑定失败:" + err.Error()})
	}

	if requestComment.Context==""{
		c.JSON(http.StatusUnprocessableEntity,gin.H{"msg":"评论不能为空"})
		return
	}
	// 获取登录用户的id
	userid, _ := c.Get("userid")
	requestComment.Userid=userid.(int)

	//获取评论用户的名字作为comment的作者
	var user models.UserInfo
	dao.DB.Where("ID=?",userid.(int)).First(&user)
	requestComment.Author=user.Username

	// 创建数据
	dao.DB.AutoMigrate(models.Comment{})
	if err := dao.DB.Create(&requestComment).Error; err != nil {
		panic(err)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"mes":"评论成功",
		"data":requestComment,
	})
	return
}

//发布Reply
func CreatReply(c *gin.Context ) {
	var requestReply models.Reply
	// 数据验证
	if err := c.ShouldBind(&requestReply); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "用户参数绑定失败:" + err.Error()})
	}
	if requestReply.Context==""{
		c.JSON(http.StatusUnprocessableEntity,gin.H{"msg":"回复不能为空"})
		return
	}
	// 获取登录用户的id
	userid, _ := c.Get("userid")
	requestReply.Userid=userid.(int)

	//获取评论用户的名字作为reply的作者
	var user models.UserInfo
	dao.DB.Where("ID=?",userid.(int)).First(&user)
	requestReply.Author=user.Username

	// 创建数据
	dao.DB.AutoMigrate(models.Reply{})
	if err := dao.DB.Create(&requestReply).Error; err != nil {
		panic(err)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"mes":"回复成功",
		"data":requestReply,
	})
	return
}
