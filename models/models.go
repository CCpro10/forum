package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

//论坛的代号
const (
	Lifecode = iota
	Technologycode
	Emotioncode
	Entertainmentcode
	Gamescode
	Fashioncode
	Literaturecode
)

//控制论坛权限
type Forumpermission struct {
	ID                uint
	Forumcode         int  `form:"forumcode" `
	Postpermission    bool `form:"postpermission" `    //论坛发帖权限
	Commentpermission bool `form:"commentpermission" ` //论坛评论权限
	Accesspermission  bool `form:"accesspermission" `  //论坛访问权限
}

//管理员列表
type Managerlist struct {
	ID        uint  `gorm:"primary_key"`
	UserID    uint `form:"userid" `       //管理者的ID
	Forumcode uint   `form:"forumcode"  `   //管理的论坛代号
}

//帖子
type Post struct {
	ID        uint   `gorm:"primary_key" `
	Forumcode int    `form:"forumcode"`
	Userid    int    `form:"userid"`
	Article   string `form:"article"`
	Context   string `form:"context"`
	Author    string `form:"author"`
	//Commentpermission bool  `form:"commentpermission"default:"ture"`
	CreatedAt time.Time
	//Replies   []Reply//所有回复
}

//帖子下面的评论
type Comment struct {
	ID        uint
	Postid    int    `form:"postid" ` //前端给来评论的id
	Userid    int    `form:"userid" ` //系统根据token识别的用户id
	Author    string `form:"author" `
	Context   string `form:"context" ` //用户输入的内容
	CreatedAt time.Time
	//用来存帖子下面回复的回复
	//Comments  []Comment
}

//帖子下面回复的回复,楼中楼
type Reply struct {
	ID        uint
	Commentid int    `form:"commentid" ` //前端给来的id
	Userid    int    `form:"userid"`    //系统根据token识别的用户id
	Author    string `form:"author"`
	Context   string `form:"context"`  //用户输入的内容
	CreatedAt time.Time
}

type UserInfo struct {
	gorm.Model
	PhoneNumber string `form:"phonenumber"` //ID让用户输入手机号,作为数据库的账户名
	Username    string `form:"username"`    //用户的呢名
	Password    string `form:"password"`
	//Forumlist   [8]int `form:"forumlist"`
	//Manager    bool  `form:"manager"`//是否为管理者
	//ManagedForum  string  `form:"managedforum"`//管理的
}

