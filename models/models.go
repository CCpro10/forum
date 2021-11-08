package models

import "github.com/jinzhu/gorm"

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

type UserInfo struct {
	gorm.Model
	PhoneNumber string  `form:"phonenumber"` //ID让用户输入手机号,作为数据库的账户名
	Username    string  `form:"username"`  //用户的呢名
	Password    string  `form:"password"`
	//Forumlist   [8]int `form:"forumlist"`
	//Manager    bool  `form:"manager"`//是否为管理者
	//ManagedForum  string  `form:"managedforum"`//管理的
}
