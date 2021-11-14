package controller

import (
	"forum/dao"
	_ "forum/dao"
	"forum/models"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm"
	"net/http"
)

//登录,接收用户名和密码,如果正确的话返回一个tokenstring
func Developerlogin(c *gin.Context) {
	// 开发者发送用户名和密码过来
	var Developer models.UserInfo
	err := c.ShouldBind(&Developer)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg":  "无效的参数"})
		return
	}

	// 校验用户名和密码是否正确
	if Developer.PhoneNumber!= "123456"|| Developer.Password!= "666666"{
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "用户名或密码错误"})
		return
	}
	// 生成Token,设置token中开发者id为999999999
	tokenString, _ := GenToken(uint(999999999))
	c.JSON(http.StatusOK, gin.H{
		"msg":  "开发者登录成功",
		"data": gin.H{"token": tokenString},
	})
	return
}

//修改管理员列表,传入userid 和forumcode
func SetManager(c *gin.Context) {
	var manager models.Managerlist
	if err := c.ShouldBind(&manager); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "无效的参数"})
		return
	}
	//检查ID是否为管理员的ID
	managerid,ok:=c.Get("userid")
	if!ok||managerid!=999999999{
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "您不是开发者,无操作权限"})
		return
	}
	//检查此用户是否存在
	var user models.UserInfo
	dao.DB.Where("ID=?", manager.UserID).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "用户ID不存在",
		})
		return
	}

	//查询管理员列表,检查是否重复创建
	var managerlist models.Managerlist
	dao.DB.Where("user_id=? AND forumcode=?", manager.UserID,manager.Forumcode).First(&managerlist)
	//不为0说明该管理员已存在
	if managerlist.ID !=0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "此用户已经设置为该论坛管理员,请勿重复操作",
		})
		return
	}

	//创建列表
	dao.DB.AutoMigrate(&models.Managerlist{})
	err :=dao.DB.Create(&manager).Error
	if err!=nil{
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg":  "管理员创建失败",
		})
	}else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "管理员创建成功",
			"data": manager,
		})
	}
}
