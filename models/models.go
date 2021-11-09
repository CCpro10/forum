package models

import (
	"github.com/jinzhu/gorm"
	"time"
)
//论坛的代号
const (
Technologycode =iota
Lifecode
Emotioncode
Entertainmentcode
Gamescode
Fashioncode
Literaturecode
)

type Forum struct {

}

//管理员列表
type Managerlist struct {
	Forumcode      uint  `from:"forunmcode"`         //管理的论坛代号
	UserID   string `from:"userphonenum"`  //管理者的手机号
}

//帖子
type Post struct {
	ID        uint   `gorm:"primary_key"`
	Forumcode int    `from:"forunmcode"`
	Userid    int    `from:"userid"`
	Article   string `from:"article"`
	Context   string `from:"context"`
	Author    string `from:"author"`
	CreatedAt time.Time
	//Replies   []Reply//所有回复
}

//帖子下面的评论
type Comment struct{
	ID 		uint
	Postid  int			 `from:"post"`
	Author  string		 `from:"author"`
	Context string		 `from:"context"`
	CreatedAt time.Time
	//用来存帖子下面回复的回复
	//Comments  []Comment
}

//帖子下面回复的回复
type  Reply struct {
	ID      uint
	Commentid int		 `from:"commentid"`
	Author  string		 `from:"author"`
	Context string		 `from:"context"`
	CreatedAt time.Time
}

type UserInfo struct {
	gorm.Model
	PhoneNumber string  `form:"phonenumber"` //ID让用户输入手机号,作为数据库的账户名
	Username    string  `form:"username"`  //用户的呢名
	Password    string  `form:"password"`
	//Forumlist   [8]int `form:"forumlist"`
	//Manager    bool  `form:"manager"`//是否为管理者
	//ManagedForum  string  `form:"managedforum"`//管理的
}

