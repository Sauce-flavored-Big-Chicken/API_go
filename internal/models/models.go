package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName     string  `json:"userName" gorm:"column:user_name"`
	NickName     string  `json:"nickName" gorm:"column:nick_name"`
	PassWord     string  `json:"passWord" gorm:"column:pass_word"`
	Phone        string  `json:"phone" gorm:"column:phone"`
	Email        string  `json:"email" gorm:"column:email"`
	Sex          string  `json:"sex" gorm:"column:sex"`
	Avatar       string  `json:"avatar" gorm:"column:avatar"`
	Status       string  `json:"status" gorm:"column:status"`
	DelFlag      string  `json:"delFlag" gorm:"column:del_flag"`
	LoginDate    string  `json:"loginDate" gorm:"column:login_date"`
	IP           string  `json:"ip" gorm:"column:ip"`
	IDCard       string  `json:"idCard" gorm:"column:id_card"`
	Address      string  `json:"address" gorm:"column:address"`
	Introduction string  `json:"introduction" gorm:"column:introduction"`
	Balance      float64 `json:"balance" gorm:"column:balance"`
	Score        int     `json:"score" gorm:"column:score"`
}

type Rotation struct {
	gorm.Model
	Title   string `json:"title" gorm:"column:title"`
	PicPath string `json:"picPath" gorm:"column:pic_path"`
	Link    string `json:"link" gorm:"column:link"`
	Type    int    `json:"type" gorm:"column:type"`
	Status  string `json:"status" gorm:"column:status"`
}

type PressCategory struct {
	gorm.Model
	Name     string `json:"name" gorm:"column:name"`
	ParentId int    `json:"parentId" gorm:"column:parent_id"`
	Sort     int    `json:"sort" gorm:"column:sort"`
	Status   string `json:"status" gorm:"column:status"`
}

type PressNews struct {
	gorm.Model
	Title       string    `json:"title" gorm:"column:title"`
	Content     string    `json:"content" gorm:"column:content;type:text"`
	ImageUrls   string    `json:"imageUrls" gorm:"column:image_urls"`
	CategoryId  int       `json:"categoryId" gorm:"column:category_id"`
	Author      string    `json:"author" gorm:"column:author"`
	Source      string    `json:"source" gorm:"column:source"`
	ViewCount   int       `json:"viewCount" gorm:"column:view_count"`
	Status      string    `json:"status" gorm:"column:status"`
	PublishDate time.Time `json:"publishDate" gorm:"column:publish_date"`
}

type Notice struct {
	gorm.Model
	Title         string    `json:"title" gorm:"column:title"`
	NoticeStatus  string    `json:"noticeStatus" gorm:"column:notice_status"`
	NoticeContent string    `json:"noticeContent" gorm:"column:notice_content;type:text"`
	PublishDate   time.Time `json:"publishDate" gorm:"column:publish_date"`
	CreateBy      string    `json:"createBy" gorm:"column:create_by"`
}

type FriendlyNeighbor struct {
	gorm.Model
	UserId     int    `json:"userId" gorm:"column:user_id"`
	NickName   string `json:"nickName" gorm:"column:nick_name"`
	UserImgUrl string `json:"userImgUrl" gorm:"column:user_img_url"`
	Content    string `json:"content" gorm:"column:content;type:text"`
	ImgUrl     string `json:"imgUrl" gorm:"column:img_url"`
	CommentNum int    `json:"commentNum" gorm:"column:comment_num"`
	LikeNum    int    `json:"likeNum" gorm:"column:like_num"`
	CreateTime string `json:"createTime" gorm:"column:create_time"`
}

type FNComment struct {
	gorm.Model
	NeighborId int    `json:"neighborId" gorm:"column:neighbor_id"`
	UserId     int    `json:"userId" gorm:"column:user_id"`
	NickName   string `json:"nickName" gorm:"column:nick_name"`
	UserImgUrl string `json:"userImgUrl" gorm:"column:user_img_url"`
	Content    string `json:"content" gorm:"column:content;type:text"`
	CreateTime string `json:"createTime" gorm:"column:create_time"`
}

type ActivityCategory struct {
	gorm.Model
	Name   string `json:"name" gorm:"column:name"`
	Status string `json:"status" gorm:"column:status"`
}

type Activity struct {
	gorm.Model
	Title        string    `json:"title" gorm:"column:title"`
	Content      string    `json:"content" gorm:"column:content;type:text"`
	PicPath      string    `json:"picPath" gorm:"column:pic_path"`
	CategoryId   int       `json:"categoryId" gorm:"column:category_id"`
	StartDate    time.Time `json:"startDate" gorm:"column:start_date"`
	EndDate      time.Time `json:"endDate" gorm:"column:end_date"`
	Address      string    `json:"address" gorm:"column:address"`
	TotalCount   int       `json:"totalCount" gorm:"column:total_count"`
	CurrentCount int       `json:"currentCount" gorm:"column:current_count"`
	Status       string    `json:"status" gorm:"column:status"`
	CreateBy     string    `json:"createBy" gorm:"column:create_by"`
	CreateTime   string    `json:"createTime" gorm:"column:create_time"`
}

type Registration struct {
	gorm.Model
	UserId        int    `json:"userId" gorm:"column:user_id"`
	UserName      string `json:"userName" gorm:"column:user_name"`
	NickName      string `json:"nickName" gorm:"column:nick_name"`
	Phone         string `json:"phone" gorm:"column:phone"`
	ActivityId    int    `json:"activityId" gorm:"column:activity_id"`
	Status        string `json:"status" gorm:"column:status"`
	CheckinStatus string `json:"checkinStatus" gorm:"column:checkin_status"`
	Comment       string `json:"comment" gorm:"column:comment"`
	CreateTime    string `json:"createTime" gorm:"column:create_time"`
}

type Comment struct {
	gorm.Model
	Type       string `json:"type" gorm:"column:type"`
	Sid        int    `json:"sid" gorm:"column:sid"`
	Content    string `json:"content" gorm:"column:content;type:text"`
	LikeNum    int    `json:"likeNum" gorm:"column:like_num"`
	ReplyNum   int    `json:"replyNum" gorm:"column:reply_num"`
	UserId     int    `json:"userId" gorm:"column:user_id"`
	NickName   string `json:"nickName" gorm:"column:nick_name"`
	UserImgUrl string `json:"userImgUrl" gorm:"column:user_img_url"`
	CreateTime string `json:"createTime" gorm:"column:create_time"`
}
