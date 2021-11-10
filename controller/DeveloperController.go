package controller

import (
	"forum/dao"
	_ "forum/dao"
	"forum/models"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm"
	"log"
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
	log.Printf("%+v",Developer)
	// 校验用户名和密码是否正确
	if Developer.PhoneNumber!= "123456"|| Developer.Password!= "666666"{
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "用户名或密码错误"})
		return
	}

	// 生成Token
	tokenString, _ := GenToken(uint(123456))
	c.JSON(http.StatusOK, gin.H{
		"msg":  "开发者登录成功",
		"data": gin.H{"token": tokenString},
	})
	return
}

//修改管理员列表
func SetManager(c *gin.Context) {
	var manager models.Managerlist
	if err := c.ShouldBind(&manager); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "无效的参数"})
		return
	}

	var user models.UserInfo
	dao.DB.Where("ID=?", manager.UserID).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "用户ID不存在",
		})
		return
	}

	dao.DB.AutoMigrate(&models.Managerlist{})
	err :=dao.DB.Create(&manager).Error//?????,创建成功无信息
	if err!=nil{
		c.JSON(http.StatusOK, gin.H{
			"msg":  "管理员创建成功",
			"data": manager,
		})
	}

}
