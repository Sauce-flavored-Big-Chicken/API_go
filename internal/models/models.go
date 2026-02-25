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
	SubTitle    string    `json:"subTitle" gorm:"column:sub_title"`
	Content     string    `json:"content" gorm:"column:content;type:text"`
	ImageUrls   string    `json:"imageUrls" gorm:"column:image_urls"`
	CategoryId  int       `json:"categoryId" gorm:"column:category_id"`
	Author      string    `json:"author" gorm:"column:author"`
	Source      string    `json:"source" gorm:"column:source"`
	ViewCount   int       `json:"viewCount" gorm:"column:view_count"`
	LikeNum     int       `json:"likeNum" gorm:"column:like_num"`
	CommentNum  int       `json:"commentNum" gorm:"column:comment_num"`
	Type        string    `json:"type" gorm:"column:type"`
	Top         string    `json:"top" gorm:"column:top"`
	Hot         string    `json:"hot" gorm:"column:hot"`
	Tags        string    `json:"tags" gorm:"column:tags"`
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
	IsTop        string    `json:"isTop" gorm:"column:is_top"`
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
	Star          int    `json:"star" gorm:"column:star"`
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

type PressLikeRecord struct {
	gorm.Model
	NewsId int `json:"newsId" gorm:"column:news_id;index:idx_press_like_user_news,unique"`
	UserId int `json:"userId" gorm:"column:user_id;index:idx_press_like_user_news,unique"`
}

type CommentLikeRecord struct {
	gorm.Model
	CommentId int `json:"commentId" gorm:"column:comment_id;index:idx_comment_like_user_comment,unique"`
	UserId    int `json:"userId" gorm:"column:user_id;index:idx_comment_like_user_comment,unique"`
}

type GreenDataCard struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Icon  string `json:"icon" gorm:"column:icon"`
	Title string `json:"title" gorm:"column:title"`
	Num   string `json:"num" gorm:"column:num"`
	Unit  string `json:"unit" gorm:"column:unit"`
	Trend string `json:"trend" gorm:"column:trend"`
	Sort  int    `json:"sort" gorm:"column:sort"`
}

type GreenQuestion struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	QuestionType string `json:"questionType" gorm:"column:question_type;index:idx_green_question_type_level"`
	Level        string `json:"level" gorm:"column:level;index:idx_green_question_type_level"`
	Question     string `json:"question" gorm:"column:question;type:text"`
	OptionA      string `json:"optionA" gorm:"column:option_a"`
	OptionB      string `json:"optionB" gorm:"column:option_b"`
	OptionC      string `json:"optionC" gorm:"column:option_c"`
	OptionD      string `json:"optionD" gorm:"column:option_d"`
	OptionE      string `json:"optionE" gorm:"column:option_e"`
	OptionF      string `json:"optionF" gorm:"column:option_f"`
	Answer       string `json:"answer" gorm:"column:answer"`
	Score        int    `json:"score" gorm:"column:score"`
	Status       string `json:"status" gorm:"column:status"`
}

type GreenPaper struct {
	gorm.Model
	UserId   int    `json:"userId" gorm:"column:user_id;index"`
	Score    string `json:"score" gorm:"column:score"`
	RawInput string `json:"rawInput" gorm:"column:raw_input;type:text"`
}

type GreenPaperAnswer struct {
	gorm.Model
	PaperId     uint   `json:"paperId" gorm:"column:paper_id;index"`
	QuestionId  uint   `json:"questionId" gorm:"column:question_id;index"`
	UserAnswer  string `json:"userAnswer" gorm:"column:user_answer"`
	RightAnswer string `json:"rightAnswer" gorm:"column:right_answer"`
	IsCorrect   string `json:"isCorrect" gorm:"column:is_correct"`
}

type GreenDataSeries struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	ListKey string `json:"listKey" gorm:"column:list_key;index:idx_green_data_series_key_sort"`
	Name    string `json:"name" gorm:"column:name"`
	Data    string `json:"data" gorm:"column:data;type:text"`
	Sort    int    `json:"sort" gorm:"column:sort;index:idx_green_data_series_key_sort"`
}
