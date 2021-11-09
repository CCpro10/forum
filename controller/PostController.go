package controller
import (
	_"fmt"
	"forum/dao"
	"github.com/gin-gonic/gin"
	_"github.com/jinzhu/gorm"
	_"golang.org/x/crypto/bcrypt"
	_"log"
	"net/http"
	"forum/models"
)

func Postcreat(c *gin.Context ) {
	var requestPost models.Post
	// 数据验证
	if err := c.ShouldBind(&requestPost); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "用户参数绑定失败:" + err.Error()})
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
	c.JSON(http.StatusOK,gin.H{"mes":"创建成功"})
	return
}
