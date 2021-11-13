package controller
import (
	"forum/dao"
	_ "forum/dao"
	"forum/models"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm"
	_"log"
	_"log"
	"net/http"
)

//设置论坛发帖权限,传入forumcode和postpermission
func SetPostPermission(c *gin.Context) {
	//接收参数
	var forumpermisson models.Forumpermission
	if err := c.ShouldBind(&forumpermisson); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "接收参数错误Err= " + err.Error()})
		return
	}
	//检验权限,先获取用户账号
	userid, ok := c.Get("userid")
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "获取用户账号失败"})
		return
	}

	//查询管理员列表,检验此人是否在管理员列表中
	var managerlist = models.Managerlist{}
	dao.DB.Where("user_id=? AND forumcode=?", userid, forumpermisson.Forumcode).First(&managerlist)
	if managerlist.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "您不是此论坛的管理员,没有操作权限"})
		return
	}

	//更新论坛权限forumpermission
	dao.DB.AutoMigrate(models.Forumpermission{})
	dao.DB.Model(&models.Forumpermission{}).Where("forumcode = ?", forumpermisson.Forumcode).Update("postpermission", forumpermisson.Postpermission)
	c.JSON(http.StatusOK, gin.H{"msg": "论坛发帖权限修改成功"})
	return

}

//设置论坛评论权限,传入forumcode和commentpermission
func SetCommentPermission(c *gin.Context) {
	//接收参数
	var forumpermission models.Forumpermission
	if err := c.ShouldBind(&forumpermission); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "接收参数错误Err= " + err.Error()})
		return
	}
	//检验权限,先获取用户账号
	userid, ok := c.Get("userid")
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "获取用户账号失败"})
		return
	}

	//查询管理员列表,检验此人是否在管理员列表中
	var managerlist = models.Managerlist{}
	dao.DB.Where("user_id=? AND forumcode=?", userid, forumpermission.Forumcode).First(&managerlist)
	if managerlist.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "您不是此论坛的管理员,没有操作权限"})
		return
	}

	//更新论坛权限forumpermission
	dao.DB.AutoMigrate(models.Forumpermission{})
	dao.DB.Model(&models.Forumpermission{}).Where("forumcode=?", forumpermission.Forumcode).Update("commentpermission", forumpermission.Commentpermission)

	c.JSON(http.StatusOK, gin.H{"msg": "论坛评论权限修改成功"})
	return
}

//设置论坛访问权限,传入forumcode和accesspermission
func SetAccessPermission(c *gin.Context) {
	//接收参数
	var forumpermisson models.Forumpermission
	if err := c.ShouldBind(&forumpermisson); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "接收参数错误Err= " + err.Error()})
		return
	}
	//检验权限,先获取用户账号
	userid, ok := c.Get("userid")
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "获取用户账号失败"})
		return
	}

	//查询管理员列表,检验此人是否在管理员列表中
	var managerlist = models.Managerlist{}
	dao.DB.Where("user_id=? AND forumcode=?", userid, forumpermisson.Forumcode).First(&managerlist)
	if managerlist.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "您不是此论坛的管理员,没有操作权限"})
		return
	}

	//更新论坛权限forumpermission
	dao.DB.AutoMigrate(models.Forumpermission{})
	dao.DB.Model(&models.Forumpermission{}).Where("forumcode = ?", forumpermisson.Forumcode).Update("accesspermission", forumpermisson.Accesspermission)
	c.JSON(http.StatusOK, gin.H{"msg": "论坛访问权限修改成功"})
	return

}