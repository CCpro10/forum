package controller

import (
	"fmt"
	"forum/dao"
    "forum/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

//发布post,接收forumcode,context,article
func CreatPost(c *gin.Context ) {
	var requestPost models.Post
	// 数据验证
	if err := c.ShouldBind(&requestPost); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "用户参数绑定失败:" + err.Error()})
	}
	//检查发帖权限,以及帖子格式
	var forum models.Forumpermission
	dao.DB.Where("forumcode = ?", requestPost.Forumcode).First(&forum)
	if forum.Postpermission == false {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "此论坛现已禁止发帖"})
		return
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


//发布comment,接收postid,context
func CreatComment(c *gin.Context ) {
	var requestComment models.Comment
	// 数据验证
	if err := c.ShouldBind(&requestComment); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "用户参数绑定失败:" + err.Error()})
	}

	//检查评论权限,以及评论格式
	var post models.Post
	dao.DB.Where("ID = ?", requestComment.Postid).First(&post)
	var forumpermission models.Forumpermission
	dao.DB.Where(" forumcode= ?", post.Forumcode).First(&forumpermission)
	log.Printf("%+v",forumpermission)
	if forumpermission.Commentpermission == false {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "此论坛的帖子现已禁止评论"})
		return
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

//发布Reply,接收commentid,context
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

//通过query传入forumcode和offsetnum
func Showposts(c *gin.Context)  {
	forumcode:=c.Query("forumcode")
	//检验论坛是否有访问权限
	var forum models.Forumpermission
	dao.DB.Where("forumcode = ?", forumcode).First(&forum)
	if forum.Accesspermission == false {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "此论坛现已禁止访问"})
		return
	}

	//offsetnum为数据库中读取的时候要跳过的数目,若offsetnum为0,则从第一个开始查
	offsetnum:=c.Query("offsetnum")
	var posts []models.Post
	err:=dao.DB.Where("forumcode=?",forumcode).Order("ID").Offset(offsetnum).Limit(6).Find(&posts).Error
	if err!=nil{
		c.JSON(http.StatusUnprocessableEntity,gin.H{"msg":fmt.Sprintf("访问失败:err=%v",err)})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"访问成功",
		"data":posts,
	})
}

//通过query传入postid和offsetnum
func Showcomments(c *gin.Context)  {

	postid:=c.Query("postid")
	//offsetnum为数据库中读取的时候要跳过的数目,若offsetnum为0,则从第一个开始查
	offsetnum:=c.Query("offsetnum")

	//通过postid找到post,再找forumcode
	var post models.Post
	dao.DB.Where("ID=?",postid).First(&post)
	forumcode:=post.Forumcode

	//检验论坛是否有访问权限
	var forum models.Forumpermission
	dao.DB.Where("forumcode = ?", forumcode).First(&forum)
	if forum.Accesspermission == false {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "此论坛现已禁止访问"})
		return
	}

	var comments []models.Comment
	//找到post的comments
	err:=dao.DB.Where("postid=?",postid).Order("ID").Offset(offsetnum).Limit(10).Find(&comments).Error
	if err!=nil{
		c.JSON(http.StatusUnprocessableEntity,gin.H{"msg":fmt.Sprintf("评论访问失败:err=%v",err)})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"访问成功",
		"data":comments,
	})
}


//通过query传入commentid和offsetnum
func Showreplies(c *gin.Context)  {

	commentid:=c.Query("commentid")
	//offsetnum为数据库中读取的时候要跳过的数目,若offsetnum为0,则从第一个开始查
	offsetnum:=c.Query("offsetnum")

	var replies []models.Reply
	//找到comment的replies
	err:=dao.DB.Where("commentid=?",commentid).Order("ID").Offset(offsetnum).Limit(10).Find(&replies).Error
	if err!=nil{
		c.JSON(http.StatusUnprocessableEntity,gin.H{"msg":fmt.Sprintf("访问失败:err=%v",err)})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"访问成功",
		"data":replies,
	})
}


