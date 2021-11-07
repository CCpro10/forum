package models

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
	//Id   int
	Username  string`form:"username" `
	Password  string `form:"password" `

}
