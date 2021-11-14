package controller
import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/dgrijalva/jwt-go"
	"forum/models"
	"forum/dao"
	"github.com/jinzhu/gorm"
)

const TokenExpireDuration = time.Hour*24
var MySecret = []byte("天王盖地虎")

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserId uint
	jwt.StandardClaims
}

//生成令牌,传入userid,生成JWT
func GenToken(UserId uint) (string, error) {
	// 创建一个我们自己的声明/请求
	c := MyClaims{
		UserId, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "my-project",                               // 签发人
			Subject:   "user token",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析tokenstring，返回一个包含信息的用户声明
func ParseToken(tokenString string) (*MyClaims, error) {

	// 通过tokenstring,请求结构,返回秘钥的一个回调函数 这三个参数,返回一个token结构体,token包含了请求结构
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}

	// 校验token,token有效则返回claims请求，nil
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	//token无效，返回错误
	return nil, errors.New("invalid token")
}

//基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"msg":  "请求头中auth为空,请先登录"})
			c.Abort()
			return
		}

		// 按空格分割,在第一个空格后分割成两部分
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "请求头中auth格式有误,请重新登录"})
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"msg":  "此Token无效或已过期,请重新登录"})
			c.Abort()
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("userid", int(mc.UserId))
		c.Next() // 后续的处理函数可以用过c.Get("userid")来获取当前请求的用户信息
	}
}


//注册路由
func Register(c *gin.Context) {
	// 前端页面填写待办事项 点击提交 会发请求到这里
	// 1. 从请求中把数据拿出来
	var user models.UserInfo
	dao.DB.AutoMigrate(models.UserInfo{})
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "用户参数绑定失败:" + err.Error()})
	}

	//2.验证id和密码的结构
	log.Println(user.PhoneNumber,len(user.PhoneNumber), user.Password, user.Username,len(user.Password))
	if len(user.PhoneNumber) != 11 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "手机号必须为11位"})
		return
	}
	if len(user.Password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "密码长度不能小于六位"})
		return
	}
	if len(user.Username) == 0 {
		user.Username = "匿名用户"
	}

	//判断在数据库中账号是否存在
	if isPhoneNumberExist(dao.DB, user.PhoneNumber) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "此手机号已被注册"})
		return
	}

	// 3. 对密码加密
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "用户密码加密失败"})
		return
	}
	user.Password=string(hashPassword)
	//4.存入注册信息
	err = dao.DB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg":"注册信息存入数据库失败:"+err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "注册成功,请你重新登录",
			"date": user,
		})
	}
}

//登录,接收用户名和密码,如果正确的话返回一个tokenstring
func Login(c *gin.Context) {
	// 用户发送用户名和密码过来
	var requestUser models.UserInfo
	var user models.UserInfo
	err := c.ShouldBind(&requestUser)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg":  "无效的参数",
		})
		return
	}
	fmt.Println(requestUser)
	// 校验用户名和密码是否正确
	if len(requestUser.PhoneNumber) != 11 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "手机号必须为11位"})
		return
	}
	if len(requestUser.Password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "密码长度不能小于六位"})
		return
	}

	//判断手机号是否存在
	dao.DB.Where("phone_number = ?", requestUser.PhoneNumber).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "用户不存在，请先注册"})
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestUser.Password)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "账户名或密码错误"})
		return
	}

	// 生成Token
	tokenString, _ := GenToken(user.ID)
	c.JSON(http.StatusOK, gin.H{
		"msg":  "登录成功",
		"data": gin.H{"token": tokenString},
	})
	return
}

func ShowHomePage(c * gin.Context)  {
		userid:=c.MustGet("userid")
	    var user models.UserInfo
		dao.DB.Where("ID=?",userid.(int)).First(&user)
		c.JSON(200,gin.H{
			"msg":"访问成功",
			"date":user,
		})

}

//传入username
func ChangeUsername(c *gin.Context)  {
	//绑定userna参数
	var qequestuser models.UserInfo
    err:=c.ShouldBind(&qequestuser)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,gin.H{"msg":"用户名参数绑定失败"})
		return
	}

	//用户名不能太短
	if len(qequestuser.Username)<4{
		c.JSON(http.StatusUnprocessableEntity,gin.H{"msg":"用户名过短"})
		return
	}

	//获取userid
	 userid,ok:=c.Get("userid")
	if !ok  {
		c.JSON(http.StatusUnprocessableEntity,gin.H{"msg":"用户ID获取失败"})
		return
	}

	//通过id修改username
	dao.DB.AutoMigrate(models.UserInfo{})
	dao.DB.Model(&models.UserInfo{}).Where("ID= ?", userid).Update("username", qequestuser.Username)
	c.JSON(http.StatusOK, gin.H{"msg": "用户名修改成功"})
	return
}




//如果手机号在数据库中已存在返回ture
func isPhoneNumberExist(db *gorm.DB, PhoneNumber string)bool{
	var user models.UserInfo
	//查询并赋值给user
	db.Where("Phone_number=?", PhoneNumber).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
