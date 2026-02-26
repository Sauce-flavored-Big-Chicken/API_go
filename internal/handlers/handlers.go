package handlers

import (
	"bytes"
	"digital-community/internal/config"
	"digital-community/internal/models"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "golang.org/x/image/webp"
	"gorm.io/gorm"
)

type Response struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
	Token string      `json:"token,omitempty"`
}

func respondList(c *gin.Context, msg string, data interface{}, total int64) {
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   msg,
		"data":  data,
		"total": total,
	})
}

type smsCodeRecord struct {
	Code      string
	ExpiresAt time.Time
}

var smsCodeStore = struct {
	sync.RWMutex
	Data map[string]smsCodeRecord
}{
	Data: map[string]smsCodeRecord{},
}

var smsRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var phonePattern = regexp.MustCompile(`^1[3-9]\d{9}$`)

func genSMSCode() string {
	return strconv.Itoa(1000 + smsRand.Intn(9000))
}

func setSMSCode(phone, code string) {
	smsCodeStore.Lock()
	defer smsCodeStore.Unlock()
	smsCodeStore.Data[phone] = smsCodeRecord{
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

func verifySMSCode(phone, code string) bool {
	smsCodeStore.RLock()
	rec, ok := smsCodeStore.Data[phone]
	smsCodeStore.RUnlock()
	if !ok {
		return false
	}
	if time.Now().After(rec.ExpiresAt) {
		smsCodeStore.Lock()
		delete(smsCodeStore.Data, phone)
		smsCodeStore.Unlock()
		return false
	}
	return rec.Code == code
}

func clearSMSCode(phone string) {
	smsCodeStore.Lock()
	defer smsCodeStore.Unlock()
	delete(smsCodeStore.Data, phone)
}

type userInfoResp struct {
	UserID       uint    `json:"userId"`
	UserName     string  `json:"userName"`
	Avatar       string  `json:"avatar"`
	NickName     string  `json:"nickName"`
	PhoneNumber  string  `json:"phoneNumber"`
	Sex          string  `json:"sex"`
	Email        string  `json:"email"`
	Balance      float64 `json:"balance"`
	Score        int     `json:"score"`
	IDCard       string  `json:"idCard"`
	Address      string  `json:"address"`
	Introduction string  `json:"introduction"`
}

func buildUserInfoResp(user models.User) userInfoResp {
	return userInfoResp{
		UserID:       user.ID,
		UserName:     user.UserName,
		Avatar:       user.Avatar,
		NickName:     user.NickName,
		PhoneNumber:  user.Phone,
		Sex:          user.Sex,
		Email:        user.Email,
		Balance:      user.Balance,
		Score:        user.Score,
		IDCard:       user.IDCard,
		Address:      user.Address,
		Introduction: user.Introduction,
	}
}

func buildPressItem(news models.PressNews) gin.H {
	publishDate := news.PublishDate.Format("2006-01-02")
	if news.PublishDate.IsZero() {
		publishDate = ""
	}
	return gin.H{
		"id":          news.ID,
		"cover":       news.ImageUrls,
		"title":       news.Title,
		"subTitle":    news.SubTitle,
		"content":     news.Content,
		"status":      news.Status,
		"publishDate": publishDate,
		"commentNum":  news.CommentNum,
		"likeNum":     news.LikeNum,
		"readNum":     news.ViewCount,
		"type":        news.Type,
		"top":         news.Top,
		"hot":         news.Hot,
		"tags":        news.Tags,
	}
}

func buildNoticeItem(notice models.Notice) gin.H {
	return gin.H{
		"id":            notice.ID,
		"noticeTitle":   notice.Title,
		"noticeStatus":  notice.NoticeStatus,
		"contentNotice": notice.NoticeContent,
		"releaseUnit":   notice.CreateBy,
		"phone":         "",
		"createTime":    notice.PublishDate.Format("2006-01-02 15:04:05"),
		"expressId":     1,
		"noticeName":    "重要通知",
	}
}

func buildNeighborItem(neighbor models.FriendlyNeighbor) gin.H {
	return gin.H{
		"id":             neighbor.ID,
		"publishName":    neighbor.NickName,
		"likeNum":        neighbor.LikeNum,
		"title":          "",
		"publishTime":    neighbor.CreateTime,
		"publishContent": neighbor.Content,
		"commentNum":     neighbor.CommentNum,
		"imgUrl":         neighbor.ImgUrl,
		"userImgUrl":     neighbor.UserImgUrl,
	}
}

func buildActivityItem(activity models.Activity) gin.H {
	isTop := activity.IsTop
	if isTop == "" {
		isTop = "0"
	}

	return gin.H{
		"id":            activity.ID,
		"category":      strconv.Itoa(activity.CategoryId),
		"title":         activity.Title,
		"picPath":       activity.PicPath,
		"startDate":     activity.StartDate.Format("2006/1/2 15:04"),
		"endDate":       activity.EndDate.Format("2006/1/2 15:04"),
		"sponsor":       activity.CreateBy,
		"content":       activity.Content,
		"signUpNum":     activity.CurrentCount,
		"maxNum":        activity.TotalCount,
		"signUpEndDate": nil,
		"isTop":         isTop,
	}
}

func parsePaging(c *gin.Context) (int, int) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return pageNum, pageSize
}

func PhoneLogin(c *gin.Context) {
	var req struct {
		Phone      string `json:"phone" binding:"required"`
		SMSCode    string `json:"smsCode"`
		LegacyCode string `json:"SMSCode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if req.SMSCode == "" {
		req.SMSCode = req.LegacyCode
	}
	if req.SMSCode == "" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if !phonePattern.MatchString(req.Phone) {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "手机号格式错误"})
		return
	}
	if !verifySMSCode(req.Phone, req.SMSCode) {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "验证码错误或已过期"})
		return
	}

	var user models.User
	if err := config.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "手机号未注册"})
		return
	}

	token, err := generateToken(int(user.ID), user.UserName, user.NickName, user.Phone)
	if err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "生成token失败"})
		return
	}
	clearSMSCode(req.Phone)

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功", Token: token})
}

func Login(c *gin.Context) {
	var req struct {
		UserName       string `json:"userName" binding:"required"`
		Password       string `json:"password"`
		LegacyPassword string `json:"passWord"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if req.Password == "" {
		req.Password = req.LegacyPassword
	}
	if req.Password == "" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	var user models.User
	result := config.DB.Where("user_name = ? AND pass_word = ?", req.UserName, req.Password).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "用户名或密码错误"})
		return
	}

	token, err := generateToken(int(user.ID), user.UserName, user.NickName, user.Phone)
	if err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功", Token: token})
}

func SMSCode(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "手机号不能为空"})
		return
	}
	if !phonePattern.MatchString(phone) {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "手机号格式错误"})
		return
	}
	var count int64
	config.DB.Model(&models.User{}).Where("phone = ?", phone).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "手机号未注册"})
		return
	}
	code := genSMSCode()
	setSMSCode(phone, code)
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "请求成功",
		Data: code,
	})
}

func Register(c *gin.Context) {
	var req struct {
		Avatar       string `json:"avatar"`
		UserName     string `json:"userName" binding:"required"`
		NickName     string `json:"nickName"`
		Password     string `json:"password"`
		PassWord     string `json:"passWord"`
		PhoneNumber  string `json:"phoneNumber"`
		Phonenumber  string `json:"phonenumber"`
		Sex          string `json:"sex" binding:"required"`
		Email        string `json:"email"`
		IDCard       string `json:"idCard"`
		Address      string `json:"address"`
		Introduction string `json:"introduction"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if req.Password == "" {
		req.Password = req.PassWord
	}
	if req.PhoneNumber == "" {
		req.PhoneNumber = req.Phonenumber
	}
	if req.Password == "" || req.PhoneNumber == "" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if !phonePattern.MatchString(req.PhoneNumber) {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "手机号格式错误"})
		return
	}
	if req.Sex != "0" && req.Sex != "1" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "性别参数错误"})
		return
	}

	var count int64
	config.DB.Model(&models.User{}).Where("user_name = ?", req.UserName).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "用户名已存在"})
		return
	}
	config.DB.Model(&models.User{}).Where("phone = ?", req.PhoneNumber).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "手机号已注册"})
		return
	}

	user := models.User{
		UserName:     req.UserName,
		PassWord:     req.Password,
		NickName:     req.NickName,
		Phone:        req.PhoneNumber,
		Sex:          req.Sex,
		Email:        req.Email,
		Avatar:       req.Avatar,
		IDCard:       req.IDCard,
		Address:      req.Address,
		Introduction: req.Introduction,
		Balance:      0,
		Score:        0,
		Status:       "0",
		DelFlag:      "0",
	}
	config.DB.Create(&user)

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func GetUserInfo(c *gin.Context) {
	userId := c.GetInt("userId")
	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "获取数据成功", Data: buildUserInfoResp(user)})
}

func UpdateUserInfo(c *gin.Context) {
	userId := c.GetInt("userId")
	var req struct {
		Avatar       string `json:"avatar"`
		NickName     string `json:"nickName" binding:"required"`
		PhoneNumber  string `json:"phoneNumber"`
		Phonenumber  string `json:"phonenumber"`
		Sex          string `json:"sex"`
		Email        string `json:"email"`
		IDCard       string `json:"idCard"`
		Address      string `json:"address"`
		Introduction string `json:"introduction"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if req.PhoneNumber == "" {
		req.PhoneNumber = req.Phonenumber
	}

	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "用户不存在"})
		return
	}

	config.DB.Model(&user).Updates(map[string]interface{}{
		"nick_name":    req.NickName,
		"phone":        req.PhoneNumber,
		"sex":          req.Sex,
		"avatar":       req.Avatar,
		"email":        req.Email,
		"id_card":      req.IDCard,
		"address":      req.Address,
		"introduction": req.Introduction,
	})

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func ResetPwd(c *gin.Context) {
	userId := c.GetInt("userId")
	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "用户不存在"})
		return
	}
	if user.PassWord != req.OldPassword {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "原密码错误"})
		return
	}

	config.DB.Model(&models.User{}).Where("id = ?", userId).Update("pass_word", req.NewPassword)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func UserList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)

	var users []models.User
	var total int64
	config.DB.Model(&models.User{}).Count(&total)

	config.DB.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&users)

	items := make([]gin.H, 0, len(users))
	for _, v := range users {
		items = append(items, gin.H{
			"id":           v.ID,
			"userName":     v.UserName,
			"nickName":     v.NickName,
			"phone":        v.Phone,
			"email":        v.Email,
			"sex":          v.Sex,
			"avatar":       v.Avatar,
			"address":      v.Address,
			"introduction": v.Introduction,
			"balance":      v.Balance,
			"score":        v.Score,
			"status":       v.Status,
			"createTime":   v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	respondList(c, "获取成功", items, total)
}

func UserCreate(c *gin.Context) {
	var req struct {
		UserName     string `json:"userName" binding:"required"`
		Password     string `json:"password" binding:"required"`
		NickName     string `json:"nickName"`
		Phone        string `json:"phone" binding:"required"`
		Sex          string `json:"sex"`
		Email        string `json:"email"`
		Avatar       string `json:"avatar"`
		Address      string `json:"address"`
		Introduction string `json:"introduction"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	var count int64
	config.DB.Model(&models.User{}).Where("user_name = ?", req.UserName).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "用户名已存在"})
		return
	}
	config.DB.Model(&models.User{}).Where("phone = ?", req.Phone).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "手机号已存在"})
		return
	}

	user := models.User{
		UserName:     req.UserName,
		PassWord:     req.Password,
		NickName:     req.NickName,
		Phone:        req.Phone,
		Sex:          req.Sex,
		Email:        req.Email,
		Avatar:       req.Avatar,
		Address:      req.Address,
		Introduction: req.Introduction,
		Balance:      0,
		Score:        0,
		Status:       "0",
		DelFlag:      "0",
	}
	config.DB.Create(&user)

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "创建成功", Data: user.ID})
}

func UserUpdate(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil || userId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	var req struct {
		NickName     string `json:"nickName"`
		Phone        string `json:"phone"`
		Sex          string `json:"sex"`
		Email        string `json:"email"`
		Avatar       string `json:"avatar"`
		Address      string `json:"address"`
		Introduction string `json:"introduction"`
		Status       string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "用户不存在"})
		return
	}

	updates := map[string]interface{}{}
	if req.NickName != "" {
		updates["nick_name"] = req.NickName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Sex != "" {
		updates["sex"] = req.Sex
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Introduction != "" {
		updates["introduction"] = req.Introduction
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	config.DB.Model(&user).Updates(updates)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "更新成功"})
}

func UserDelete(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil || userId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "用户不存在"})
		return
	}

	config.DB.Delete(&user)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "退出成功"})
}

func RotationList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)
	rtype := c.Query("type")
	if rtype == "" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	var rotations []models.Rotation
	query := config.DB.Model(&models.Rotation{}).Where("type = ?", rtype)

	var total int64
	query.Count(&total)

	query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&rotations)

	items := make([]gin.H, 0, len(rotations))
	for _, r := range rotations {
		items = append(items, gin.H{
			"id":       r.ID,
			"advTitle": r.Title,
			"advImg":   r.PicPath,
			"type":     strconv.Itoa(r.Type),
		})
	}
	respondList(c, "请求成功", items, total)
}

func RotationCreate(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		PicPath string `json:"picPath"`
		Link    string `json:"link"`
		Type    int    `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	rotation := models.Rotation{
		Title:   req.Title,
		PicPath: req.PicPath,
		Link:    req.Link,
		Type:    req.Type,
		Status:  "0",
	}
	config.DB.Create(&rotation)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "创建成功", Data: rotation.ID})
}

func RotationUpdate(c *gin.Context) {
	id := c.Param("id")
	rotationId, err := strconv.Atoi(id)
	if err != nil || rotationId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	var req struct {
		Title   string `json:"title"`
		PicPath string `json:"picPath"`
		Link    string `json:"link"`
		Type    int    `json:"type"`
		Status  string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.PicPath != "" {
		updates["pic_path"] = req.PicPath
	}
	if req.Link != "" {
		updates["link"] = req.Link
	}
	if req.Type > 0 {
		updates["type"] = req.Type
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	result := config.DB.Model(&models.Rotation{}).Where("id = ?", rotationId).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "更新失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "轮播图不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "更新成功"})
}

func RotationDelete(c *gin.Context) {
	id := c.Param("id")
	rotationId, err := strconv.Atoi(id)
	if err != nil || rotationId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	result := config.DB.Delete(&models.Rotation{}, rotationId)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "删除失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "轮播图不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func PressCategoryList(c *gin.Context) {
	var categories []models.PressCategory
	var total int64
	config.DB.Model(&models.PressCategory{}).Count(&total)
	config.DB.Find(&categories)
	items := make([]gin.H, 0, len(categories))
	for _, v := range categories {
		items = append(items, gin.H{"id": v.ID, "name": v.Name, "sort": v.Sort, "appType": "smart_city"})
	}
	respondList(c, "查询成功", items, total)
}

func PressCategoryCreate(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Sort   int    `json:"sort"`
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	category := models.PressCategory{
		Name:   req.Name,
		Sort:   req.Sort,
		Status: req.Status,
	}
	config.DB.Create(&category)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "创建成功", Data: category.ID})
}

func PressCategoryUpdate(c *gin.Context) {
	id := c.Param("id")
	catId, err := strconv.Atoi(id)
	if err != nil || catId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	var req struct {
		Name   string `json:"name"`
		Sort   int    `json:"sort"`
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Sort > 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	result := config.DB.Model(&models.PressCategory{}).Where("id = ?", catId).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "更新失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "分类不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "更新成功"})
}

func PressCategoryDelete(c *gin.Context) {
	id := c.Param("id")
	catId, err := strconv.Atoi(id)
	if err != nil || catId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	result := config.DB.Delete(&models.PressCategory{}, catId)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "删除失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "分类不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func PressNewsList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)

	var newsList []models.PressNews
	var total int64
	config.DB.Model(&models.PressNews{}).Count(&total)

	config.DB.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&newsList)
	items := make([]gin.H, 0, len(newsList))
	for _, v := range newsList {
		items = append(items, buildPressItem(v))
	}
	respondList(c, "查询成功", items, total)
}

func PressCategoryNewsList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	var newsList []models.PressNews
	query := config.DB.Model(&models.PressNews{})
	query = query.Where("category_id = ?", id)

	var total int64
	query.Count(&total)

	query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&newsList)
	items := make([]gin.H, 0, len(newsList))
	for _, v := range newsList {
		items = append(items, buildPressItem(v))
	}
	respondList(c, "查询成功", items, total)
}

func PressNewsDetail(c *gin.Context) {
	id := c.Param("id")
	var news models.PressNews
	if err := config.DB.First(&news, id).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "新闻不存在"})
		return
	}
	config.DB.Model(&news).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	news.ViewCount += 1
	item := buildPressItem(news)
	item["appType"] = "community"
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "请求成功", Data: item})
}

func PressNewsCreate(c *gin.Context) {
	var req struct {
		Title      string `json:"title" binding:"required"`
		SubTitle   string `json:"subTitle"`
		Content    string `json:"content" binding:"required"`
		CategoryId int    `json:"categoryId" binding:"required"`
		Type       string `json:"type"`
		ImageUrls  string `json:"imageUrls"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	news := models.PressNews{
		Title:       req.Title,
		SubTitle:    req.SubTitle,
		Content:     req.Content,
		CategoryId:  req.CategoryId,
		Type:        req.Type,
		ImageUrls:   req.ImageUrls,
		Status:      "0",
		PublishDate: time.Now(),
	}
	config.DB.Create(&news)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "创建成功", Data: news.ID})
}

func PressNewsUpdate(c *gin.Context) {
	id := c.Param("id")
	newsId, err := strconv.Atoi(id)
	if err != nil || newsId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	var req struct {
		Title      string `json:"title"`
		SubTitle   string `json:"subTitle"`
		Content    string `json:"content"`
		CategoryId int    `json:"categoryId"`
		Type       string `json:"type"`
		ImageUrls  string `json:"imageUrls"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.SubTitle != "" {
		updates["sub_title"] = req.SubTitle
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.CategoryId > 0 {
		updates["category_id"] = req.CategoryId
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.ImageUrls != "" {
		updates["image_urls"] = req.ImageUrls
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	result := config.DB.Model(&models.PressNews{}).Where("id = ?", newsId).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "更新失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "新闻不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "更新成功"})
}

func PressNewsDelete(c *gin.Context) {
	id := c.Param("id")
	newsId, err := strconv.Atoi(id)
	if err != nil || newsId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	result := config.DB.Delete(&models.PressNews{}, newsId)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "删除失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "新闻不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func PressLike(c *gin.Context) {
	newsID, err := strconv.Atoi(c.Param("id"))
	if err != nil || newsID <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	userId := c.GetInt("userId")

	var news models.PressNews
	if err := config.DB.First(&news, newsID).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "新闻不存在"})
		return
	}

	like := models.PressLikeRecord{NewsId: newsID, UserId: userId}
	tx := config.DB.Where("news_id = ? AND user_id = ?", newsID, userId).FirstOrCreate(&like)
	if tx.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "操作失败"})
		return
	}
	if tx.RowsAffected > 0 {
		config.DB.Model(&models.PressNews{}).Where("id = ?", newsID).UpdateColumn("like_num", gorm.Expr("like_num + ?", 1))
		c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
	} else {
		c.JSON(http.StatusOK, Response{Code: 200, Msg: "已经点赞过了"})
	}
}

func PressComment(c *gin.Context) {
	var req struct {
		NewsID   string `json:"newsId" binding:"required"`
		Content  string `json:"content" binding:"required"`
		UserName string `json:"userName" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	userId := c.GetInt("userId")
	newsId, err := strconv.Atoi(req.NewsID)
	if err != nil || newsId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	comment := models.Comment{
		Sid:        newsId,
		Content:    req.Content,
		UserId:     userId,
		NickName:   req.UserName,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	var news models.PressNews
	if err := config.DB.First(&news, newsId).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "新闻不存在"})
		return
	}
	config.DB.Create(&comment)
	config.DB.Model(&models.PressNews{}).Where("id = ?", newsId).UpdateColumn("comment_num", gorm.Expr("comment_num + ?", 1))

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func CommentList(c *gin.Context) {
	id := c.Param("id")
	pageNum, pageSize := parsePaging(c)

	var comments []models.Comment
	var total int64
	config.DB.Model(&models.Comment{}).Where("sid = ?", id).Count(&total)

	config.DB.Where("sid = ?", id).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&comments)
	items := make([]gin.H, 0, len(comments))
	for _, v := range comments {
		items = append(items, gin.H{
			"id":          v.ID,
			"content":     v.Content,
			"commentDate": v.CreateTime,
			"userId":      v.UserId,
			"newsId":      v.Sid,
			"likeNum":     v.LikeNum,
		})
	}
	respondList(c, "获取数据成功", items, total)
}

func CommentLike(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil || commentID <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	userId := c.GetInt("userId")

	var comment models.Comment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "评论不存在"})
		return
	}

	like := models.CommentLikeRecord{CommentId: commentID, UserId: userId}
	tx := config.DB.Where("comment_id = ? AND user_id = ?", commentID, userId).FirstOrCreate(&like)
	if tx.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "操作失败"})
		return
	}
	if tx.RowsAffected > 0 {
		config.DB.Model(&models.Comment{}).Where("id = ?", commentID).UpdateColumn("like_num", gorm.Expr("like_num + ?", 1))
		c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
	} else {
		c.JSON(http.StatusOK, Response{Code: 200, Msg: "已经点赞过了"})
	}
}

func thumbnailURLForImage(urlPath string) string {
	cleanPath := filepath.ToSlash(filepath.Clean(urlPath))
	if !strings.HasPrefix(cleanPath, "/profile/upload/") {
		return ""
	}
	if strings.HasPrefix(cleanPath, "/profile/upload/file/") {
		return ""
	}
	if strings.HasPrefix(cleanPath, "/profile/upload/thumb/") {
		return cleanPath
	}
	rel := strings.TrimPrefix(cleanPath, "/profile/upload/")
	relNoExt := strings.TrimSuffix(rel, filepath.Ext(rel))
	return "/profile/upload/thumb/" + relNoExt + ".jpg"
}

func isImageExt(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".gif" || ext == ".svg"
}

func resizeImageNearest(src image.Image, targetW, targetH int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	for y := 0; y < targetH; y++ {
		srcY := srcBounds.Min.Y + (y * srcH / targetH)
		for x := 0; x < targetW; x++ {
			srcX := srcBounds.Min.X + (x * srcW / targetW)
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}

	return dst
}

func flattenToWhite(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(dst, dst.Bounds(), src, bounds.Min, draw.Over)
	return dst
}

func encodeJPEGUnderLimit(img image.Image, maxBytes int) ([]byte, error) {
	qualities := []int{70, 60, 50, 40, 32}
	var best []byte

	for _, q := range qualities {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: q}); err != nil {
			return nil, err
		}

		if len(best) == 0 || buf.Len() < len(best) {
			best = append([]byte(nil), buf.Bytes()...)
		}
		if buf.Len() <= maxBytes {
			return append([]byte(nil), buf.Bytes()...), nil
		}
	}

	return best, nil
}

func ensureThumbnail(urlPath string) string {
	thumbURL := thumbnailURLForImage(urlPath)
	if thumbURL == "" {
		return urlPath
	}

	sourcePath := "." + urlPath
	thumbPath := "." + thumbURL

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil || sourceInfo.IsDir() {
		return urlPath
	}

	if thumbInfo, err := os.Stat(thumbPath); err == nil && !thumbInfo.IsDir() && !thumbInfo.ModTime().Before(sourceInfo.ModTime()) {
		return thumbURL
	}

	if err := os.MkdirAll(filepath.Dir(thumbPath), 0755); err != nil {
		return urlPath
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return urlPath
	}
	defer sourceFile.Close()

	img, _, err := image.Decode(sourceFile)
	if err != nil {
		return urlPath
	}

	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	if srcW <= 0 || srcH <= 0 {
		return urlPath
	}

	const maxSide = 140
	const targetBytes = 28 * 1024
	targetW, targetH := srcW, srcH
	if srcW > maxSide || srcH > maxSide {
		if srcW >= srcH {
			targetW = maxSide
			targetH = srcH * maxSide / srcW
		} else {
			targetH = maxSide
			targetW = srcW * maxSide / srcH
		}
	}
	if targetW < 1 {
		targetW = 1
	}
	if targetH < 1 {
		targetH = 1
	}

	thumbImage := flattenToWhite(resizeImageNearest(img, targetW, targetH))
	thumbData, err := encodeJPEGUnderLimit(thumbImage, targetBytes)
	if err != nil || len(thumbData) == 0 {
		return urlPath
	}

	if err := os.WriteFile(thumbPath, thumbData, 0644); err != nil {
		_ = os.Remove(thumbPath)
		return urlPath
	}

	return thumbURL
}

func StartThumbnailWarmup() {
	go func() {
		uploadDir := "./profile/upload"
		_ = filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				if filepath.ToSlash(path) == filepath.ToSlash(filepath.Join(uploadDir, "thumb")) {
					return filepath.SkipDir
				}
				if filepath.ToSlash(path) == filepath.ToSlash(filepath.Join(uploadDir, "file")) {
					return filepath.SkipDir
				}
				return nil
			}

			ext := strings.ToLower(filepath.Ext(path))
			if !isImageExt(ext) {
				return nil
			}

			relPath := strings.TrimPrefix(filepath.ToSlash(path), ".")
			if !strings.HasPrefix(relPath, "/") {
				relPath = "/" + relPath
			}
			_ = ensureThumbnail(relPath)
			return nil
		})
	}()
}

func ImageList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)
	if pageSize > 60 {
		pageSize = 60
	}

	uploadDir := "./profile/upload"
	type imageMeta struct {
		name    string
		url     string
		size    int64
		modTime time.Time
	}

	metas := make([]imageMeta, 0, 64)

	filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if filepath.ToSlash(path) == filepath.ToSlash(filepath.Join(uploadDir, "thumb")) {
				return filepath.SkipDir
			}
			if filepath.ToSlash(path) == filepath.ToSlash(filepath.Join(uploadDir, "file")) {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if isImageExt(ext) {
			relPath := strings.TrimPrefix(filepath.ToSlash(path), ".")
			if !strings.HasPrefix(relPath, "/") {
				relPath = "/" + relPath
			}
			metas = append(metas, imageMeta{name: info.Name(), url: relPath, size: info.Size(), modTime: info.ModTime()})
		}
		return nil
	})

	sort.Slice(metas, func(i, j int) bool {
		if metas[i].modTime.Equal(metas[j].modTime) {
			return metas[i].url > metas[j].url
		}
		return metas[i].modTime.After(metas[j].modTime)
	})

	total := len(metas)
	start := (pageNum - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	images := make([]gin.H, 0, end-start)
	for _, m := range metas[start:end] {
		thumbURL := thumbnailURLForImage(m.url)
		if thumbURL == "" {
			thumbURL = m.url
		} else if _, err := os.Stat("." + thumbURL); err == nil {
			// use existing thumbnail
		} else {
			thumbURL = ensureThumbnail(m.url)
		}

		images = append(images, gin.H{
			"name":     m.name,
			"url":      m.url,
			"thumbUrl": thumbURL,
			"size":     m.size,
			"created":  m.modTime.Format("2006-01-02 15:04:05"),
		})
	}

	respondList(c, "获取成功", images, int64(total))
}

func ImageDelete(c *gin.Context) {
	urlPath := c.Query("url")
	cleanedURL, targetPath, err := resolveDeletePath(urlPath, "image")
	if err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	info, err := os.Stat(targetPath)
	if err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "文件不存在"})
		return
	}
	if info.IsDir() {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if err := os.Remove(targetPath); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "删除失败"})
		return
	}

	thumbURL := thumbnailURLForImage(cleanedURL)
	if thumbURL != "" {
		_ = os.Remove("." + thumbURL)
	}

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "上传文件失败"})
		return
	}

	name := strings.ReplaceAll(file.Filename, " ", "_")
	ext := strings.ToLower(filepath.Ext(name))
	baseDir := "image"
	if !isImageExt(ext) {
		baseDir = "file"
	}
	filename := "/profile/upload/" + baseDir + "/" + time.Now().Format("20060102150405") + "_" + name
	targetPath := "." + filename
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "上传文件失败"})
		return
	}
	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "上传文件失败"})
		return
	}

	if baseDir == "image" {
		_ = ensureThumbnail(filename)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "操作成功",
		"fileName": file.Filename,
		"url":      filename,
	})
}

func FileList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)
	if pageSize > 60 {
		pageSize = 60
	}

	uploadDir := "./profile/upload/file"
	type fileMeta struct {
		name    string
		url     string
		size    int64
		modTime time.Time
	}

	metas := make([]fileMeta, 0, 64)

	filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		relPath := strings.TrimPrefix(filepath.ToSlash(path), ".")
		if !strings.HasPrefix(relPath, "/") {
			relPath = "/" + relPath
		}
		metas = append(metas, fileMeta{name: info.Name(), url: relPath, size: info.Size(), modTime: info.ModTime()})
		return nil
	})

	sort.Slice(metas, func(i, j int) bool {
		if metas[i].modTime.Equal(metas[j].modTime) {
			return metas[i].url > metas[j].url
		}
		return metas[i].modTime.After(metas[j].modTime)
	})

	total := len(metas)
	start := (pageNum - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	files := make([]gin.H, 0, end-start)
	for _, m := range metas[start:end] {
		files = append(files, gin.H{
			"name":    m.name,
			"url":     m.url,
			"size":    m.size,
			"created": m.modTime.Format("2006-01-02 15:04:05"),
		})
	}

	respondList(c, "获取成功", files, int64(total))
}

func FileDelete(c *gin.Context) {
	urlPath := c.Query("url")
	_, targetPath, err := resolveDeletePath(urlPath, "file")
	if err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	info, err := os.Stat(targetPath)
	if err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "文件不存在"})
		return
	}
	if info.IsDir() {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if err := os.Remove(targetPath); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "删除失败"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func resolveDeletePath(urlPath, kind string) (string, string, error) {
	trimmed := strings.TrimSpace(urlPath)
	if trimmed == "" {
		return "", "", fmt.Errorf("empty path")
	}

	cleanedURL := filepath.ToSlash(filepath.Clean("/" + strings.TrimPrefix(trimmed, "/")))
	basePrefix := "/profile/upload/" + kind + "/"
	if !strings.HasPrefix(cleanedURL, basePrefix) {
		return "", "", fmt.Errorf("invalid path")
	}

	baseDir := "." + strings.TrimSuffix(basePrefix, "/")
	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", "", err
	}
	targetAbs, err := filepath.Abs("." + cleanedURL)
	if err != nil {
		return "", "", err
	}
	if targetAbs == baseAbs || !strings.HasPrefix(targetAbs, baseAbs+string(os.PathSeparator)) {
		return "", "", fmt.Errorf("path traversal")
	}

	return cleanedURL, targetAbs, nil
}

func NoticeList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)
	noticeStatus := c.Query("noticeStatus")

	var notices []models.Notice
	query := config.DB.Model(&models.Notice{})
	if noticeStatus != "" {
		query = query.Where("notice_status = ?", noticeStatus)
	}

	var total int64
	query.Count(&total)

	query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("publish_date DESC").Find(&notices)
	items := make([]gin.H, 0, len(notices))
	for _, v := range notices {
		items = append(items, buildNoticeItem(v))
	}
	respondList(c, "请求成功", items, total)
}

func NoticeDetail(c *gin.Context) {
	id := c.Param("id")
	var notice models.Notice
	if err := config.DB.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "公告不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "请求成功", Data: buildNoticeItem(notice)})
}

func NoticeCreate(c *gin.Context) {
	var req struct {
		Title         string `json:"title" binding:"required"`
		NoticeContent string `json:"noticeContent" binding:"required"`
		NoticeStatus  string `json:"noticeStatus"`
		CreateBy      string `json:"createBy"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	notice := models.Notice{
		Title:         req.Title,
		NoticeContent: req.NoticeContent,
		NoticeStatus:  req.NoticeStatus,
		CreateBy:      req.CreateBy,
		PublishDate:   time.Now(),
	}
	config.DB.Create(&notice)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "创建成功", Data: notice.ID})
}

func NoticeUpdate(c *gin.Context) {
	id := c.Param("id")
	noticeId, err := strconv.Atoi(id)
	if err != nil || noticeId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	var req struct {
		Title         string `json:"title"`
		NoticeContent string `json:"noticeContent"`
		NoticeStatus  string `json:"noticeStatus"`
		CreateBy      string `json:"createBy"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.NoticeContent != "" {
		updates["notice_content"] = req.NoticeContent
	}
	if req.NoticeStatus != "" {
		updates["notice_status"] = req.NoticeStatus
	}
	if req.CreateBy != "" {
		updates["create_by"] = req.CreateBy
	}
	result := config.DB.Model(&models.Notice{}).Where("id = ?", noticeId).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "更新失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "公告不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "更新成功"})
}

func NoticeDelete(c *gin.Context) {
	id := c.Param("id")
	noticeId, err := strconv.Atoi(id)
	if err != nil || noticeId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	result := config.DB.Delete(&models.Notice{}, noticeId)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "删除失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "公告不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func ReadNotice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	result := config.DB.Model(&models.Notice{}).Where("id = ?", id).Update("notice_status", "1")
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "操作失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "公告不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func FriendlyNeighborList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)

	var neighbors []models.FriendlyNeighbor
	var total int64
	config.DB.Model(&models.FriendlyNeighbor{}).Count(&total)

	config.DB.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("create_time DESC").Find(&neighbors)
	items := make([]gin.H, 0, len(neighbors))
	for _, v := range neighbors {
		items = append(items, buildNeighborItem(v))
	}
	respondList(c, "请求成功", items, total)
}

func FriendlyNeighborAddComment(c *gin.Context) {
	var req struct {
		NeighborhoodID int    `json:"neighborhoodId" binding:"required"`
		Content        string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	userId := c.GetInt("userId")
	nickName := c.GetString("nickName")
	if nickName == "" {
		nickName = "匿名用户"
	}

	comment := models.FNComment{
		NeighborId: req.NeighborhoodID,
		UserId:     userId,
		NickName:   nickName,
		Content:    req.Content,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	config.DB.Create(&comment)

	config.DB.Model(&models.FriendlyNeighbor{}).Where("id = ?", req.NeighborhoodID).UpdateColumn("comment_num", gorm.Expr("comment_num + ?", 1))

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func FriendlyNeighborDetail(c *gin.Context) {
	id := c.Param("id")
	var neighbor models.FriendlyNeighbor
	if err := config.DB.First(&neighbor, id).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "记录不存在"})
		return
	}

	var comments []models.FNComment
	config.DB.Where("neighbor_id = ?", id).Find(&comments)

	commentItems := make([]gin.H, 0, len(comments))
	for _, v := range comments {
		commentItems = append(commentItems, gin.H{
			"id":             v.ID,
			"userName":       v.NickName,
			"userId":         v.UserId,
			"avatar":         v.UserImgUrl,
			"content":        v.Content,
			"likeNum":        0,
			"publishTime":    v.CreateTime,
			"neighborhoodId": v.NeighborId,
		})
	}
	item := buildNeighborItem(neighbor)
	item["userComment"] = commentItems
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "请求成功", Data: item})
}

func FriendlyNeighborCreate(c *gin.Context) {
	var req struct {
		Content    string `json:"content" binding:"required"`
		ImgUrl     string `json:"imgUrl"`
		UserId     int    `json:"userId"`
		NickName   string `json:"nickName"`
		UserImgUrl string `json:"userImgUrl"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	neighbor := models.FriendlyNeighbor{
		Content:    req.Content,
		ImgUrl:     req.ImgUrl,
		UserId:     req.UserId,
		NickName:   req.NickName,
		UserImgUrl: req.UserImgUrl,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	config.DB.Create(&neighbor)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "创建成功", Data: neighbor.ID})
}

func FriendlyNeighborUpdate(c *gin.Context) {
	id := c.Param("id")
	neighborId, err := strconv.Atoi(id)
	if err != nil || neighborId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	var req struct {
		Content    string `json:"content"`
		NickName   string `json:"nickName"`
		ImgUrl     string `json:"imgUrl"`
		UserImgUrl string `json:"userImgUrl"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	updates := map[string]interface{}{}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.NickName != "" {
		updates["nick_name"] = req.NickName
	}
	if req.ImgUrl != "" {
		updates["img_url"] = req.ImgUrl
	}
	if req.UserImgUrl != "" {
		updates["user_img_url"] = req.UserImgUrl
	}
	result := config.DB.Model(&models.FriendlyNeighbor{}).Where("id = ?", neighborId).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "更新失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "帖子不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "更新成功"})
}

func FriendlyNeighborDelete(c *gin.Context) {
	id := c.Param("id")
	neighborId, err := strconv.Atoi(id)
	if err != nil || neighborId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	result := config.DB.Delete(&models.FriendlyNeighbor{}, neighborId)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "删除失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "帖子不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func ActivityTopList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)

	var activities []models.Activity
	var total int64
	query := config.DB.Model(&models.Activity{}).Where("is_top = ?", "1")
	query.Count(&total)
	if total == 0 {
		query = config.DB.Model(&models.Activity{})
		query.Count(&total)
	}

	query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("start_date DESC").Find(&activities)
	items := make([]gin.H, 0, len(activities))
	for _, v := range activities {
		items = append(items, buildActivityItem(v))
	}
	respondList(c, "请求成功", items, total)
}

func ActivityList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)

	var activities []models.Activity
	var total int64
	config.DB.Model(&models.Activity{}).Count(&total)

	config.DB.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&activities)
	items := make([]gin.H, 0, len(activities))
	for _, v := range activities {
		items = append(items, buildActivityItem(v))
	}
	respondList(c, "请求成功", items, total)
}

func ActivitySearch(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)
	var req struct {
		Words string `json:"words" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	var activities []models.Activity
	query := config.DB.Model(&models.Activity{})
	if req.Words != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+req.Words+"%", "%"+req.Words+"%")
	}

	var total int64
	query.Count(&total)

	query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&activities)
	items := make([]gin.H, 0, len(activities))
	for _, v := range activities {
		items = append(items, buildActivityItem(v))
	}
	respondList(c, "请求成功", items, total)
}

func ActivityDetail(c *gin.Context) {
	id := c.Param("id")
	var activity models.Activity
	if err := config.DB.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "活动不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "请求成功", Data: buildActivityItem(activity)})
}

func ActivityCategoryList(c *gin.Context) {
	id := c.Param("id")
	pageNum, pageSize := parsePaging(c)

	var activities []models.Activity
	query := config.DB.Model(&models.Activity{}).Where("category_id = ?", id)

	var total int64
	query.Count(&total)

	query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&activities)
	items := make([]gin.H, 0, len(activities))
	for _, v := range activities {
		items = append(items, buildActivityItem(v))
	}
	respondList(c, "请求成功", items, total)
}

func ActivityCreate(c *gin.Context) {
	var req struct {
		Title      string    `json:"title" binding:"required"`
		Content    string    `json:"content" binding:"required"`
		PicPath    string    `json:"picPath"`
		CategoryId int       `json:"categoryId" binding:"required"`
		StartDate  time.Time `json:"startDate"`
		EndDate    time.Time `json:"endDate"`
		Address    string    `json:"address"`
		TotalCount int       `json:"totalCount"`
		IsTop      string    `json:"isTop"`
		CreateBy   string    `json:"createBy"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	activity := models.Activity{
		Title:        req.Title,
		Content:      req.Content,
		PicPath:      req.PicPath,
		CategoryId:   req.CategoryId,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Address:      req.Address,
		TotalCount:   req.TotalCount,
		CurrentCount: 0,
		IsTop:        req.IsTop,
		Status:       "0",
		CreateBy:     req.CreateBy,
	}
	config.DB.Create(&activity)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "创建成功", Data: activity.ID})
}

func ActivityUpdate(c *gin.Context) {
	id := c.Param("id")
	activityId, err := strconv.Atoi(id)
	if err != nil || activityId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	var req struct {
		Title      string    `json:"title"`
		Content    string    `json:"content"`
		PicPath    string    `json:"picPath"`
		CategoryId int       `json:"categoryId"`
		StartDate  time.Time `json:"startDate"`
		EndDate    time.Time `json:"endDate"`
		Address    string    `json:"address"`
		TotalCount int       `json:"totalCount"`
		IsTop      string    `json:"isTop"`
		Status     string    `json:"status"`
		CreateBy   string    `json:"createBy"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.PicPath != "" {
		updates["pic_path"] = req.PicPath
	}
	if req.CategoryId > 0 {
		updates["category_id"] = req.CategoryId
	}
	if !req.StartDate.IsZero() {
		updates["start_date"] = req.StartDate
	}
	if !req.EndDate.IsZero() {
		updates["end_date"] = req.EndDate
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.TotalCount > 0 {
		updates["total_count"] = req.TotalCount
	}
	if req.IsTop != "" {
		updates["is_top"] = req.IsTop
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.CreateBy != "" {
		updates["create_by"] = req.CreateBy
	}
	result := config.DB.Model(&models.Activity{}).Where("id = ?", activityId).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "更新失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "活动不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "更新成功"})
}

func ActivityDelete(c *gin.Context) {
	id := c.Param("id")
	activityId, err := strconv.Atoi(id)
	if err != nil || activityId <= 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	result := config.DB.Delete(&models.Activity{}, activityId)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "删除失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "活动不存在"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func RegistrationList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)
	activityID := c.Query("activityId")
	userID := c.Query("userId")

	var registrations []models.Registration
	query := config.DB.Model(&models.Registration{})
	if activityID != "" {
		query = query.Where("activity_id = ?", activityID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	query.Count(&total)
	query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&registrations)

	items := make([]gin.H, 0, len(registrations))
	for _, v := range registrations {
		items = append(items, gin.H{
			"id":            v.ID,
			"userId":        v.UserId,
			"userName":      v.UserName,
			"nickName":      v.NickName,
			"phone":         v.Phone,
			"activityId":    v.ActivityId,
			"status":        v.Status,
			"checkinStatus": v.CheckinStatus,
			"comment":       v.Comment,
			"star":          v.Star,
			"createTime":    v.CreateTime,
		})
	}

	respondList(c, "请求成功", items, total)
}

func Registration(c *gin.Context) {
	var req struct {
		ActivityId int `json:"activityId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	userId := c.GetInt("userId")
	userName := c.GetString("userName")
	nickName := c.GetString("nickName")

	var count int64
	config.DB.Model(&models.Registration{}).Where("user_id = ? AND activity_id = ?", userId, req.ActivityId).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "已报名"})
		return
	}

	var activity models.Activity
	if err := config.DB.First(&activity, req.ActivityId).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "活动不存在"})
		return
	}
	if activity.TotalCount > 0 && activity.CurrentCount >= activity.TotalCount {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "活动报名人数已满"})
		return
	}

	registration := models.Registration{
		UserId:     userId,
		UserName:   userName,
		NickName:   nickName,
		Phone:      c.GetString("phone"),
		ActivityId: req.ActivityId,
		Status:     "0",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	config.DB.Create(&registration)

	config.DB.Model(&models.Activity{}).Where("id = ?", req.ActivityId).UpdateColumn("current_count", gorm.Expr("current_count + ?", 1))

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func Checkin(c *gin.Context) {
	activityId := c.Param("id")
	userId := c.GetInt("userId")
	result := config.DB.Model(&models.Registration{}).Where("activity_id = ? AND user_id = ?", activityId, userId).Update("checkin_status", "1")
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "操作失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "未找到报名记录"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func RegistrationComment(c *gin.Context) {
	activityId := c.Param("id")
	userId := c.GetInt("userId")
	var req struct {
		Evaluate string `json:"evaluate" binding:"required"`
		Star     int    `json:"star" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if req.Star < 1 || req.Star > 5 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "评分参数错误"})
		return
	}

	result := config.DB.Model(&models.Registration{}).Where("activity_id = ? AND user_id = ?", activityId, userId).Updates(map[string]interface{}{
		"comment": req.Evaluate,
		"star":    req.Star,
	})
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "操作失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "未找到报名记录"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func CommonDataCard(c *gin.Context) {
	var cards []models.GreenDataCard
	if err := config.DB.Order("sort asc, id asc").Find(&cards).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "查询失败"})
		return
	}
	items := make([]gin.H, 0, len(cards))
	for _, item := range cards {
		items = append(items, gin.H{
			"id":    item.ID,
			"icon":  item.Icon,
			"title": item.Title,
			"num":   item.Num,
			"unit":  item.Unit,
			"trend": item.Trend,
		})
	}
	respondList(c, "请求成功", items, int64(len(items)))
}

func QuestionQuestionList(c *gin.Context) {
	questionType := c.Param("id")
	level := c.Param("level")
	if questionType != "1" && questionType != "4" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "题型参数错误"})
		return
	}
	if level != "1" && level != "2" && level != "3" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "难度参数错误"})
		return
	}

	var total int64
	query := config.DB.Model(&models.GreenQuestion{}).
		Where("question_type = ? AND level = ? AND status = ?", questionType, level, "0")
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "查询失败"})
		return
	}

	if total == 0 {
		respondList(c, "请求成功", []gin.H{}, 0)
		return
	}

	var list []models.GreenQuestion
	if err := query.Order("RANDOM()").Limit(5).Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "查询失败"})
		return
	}

	result := make([]gin.H, 0, len(list))
	for _, q := range list {
		result = append(result, gin.H{
			"id":           q.ID,
			"questionType": q.QuestionType,
			"question":     q.Question,
			"optionA":      q.OptionA,
			"optionB":      q.OptionB,
			"optionC":      q.OptionC,
			"optionD":      q.OptionD,
			"optionE":      q.OptionE,
			"optionF":      q.OptionF,
			"answer":       q.Answer,
			"score":        q.Score,
		})
	}

	respondList(c, "请求成功", result, total)
}

func QuestionSavePaper(c *gin.Context) {
	var req struct {
		Score  interface{} `json:"score"`
		Answer []struct {
			Qid    int    `json:"qid"`
			Answer string `json:"answer"`
		} `json:"answer"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if req.Score == nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "分数不能为空"})
		return
	}
	if len(req.Answer) == 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "答案不能为空"})
		return
	}

	rawAnswers, _ := json.Marshal(req.Answer)
	paper := models.GreenPaper{
		UserId:   c.GetInt("userId"),
		Score:    strings.TrimSpace(fmt.Sprint(req.Score)),
		RawInput: string(rawAnswers),
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "提交失败"})
		return
	}
	if err := tx.Create(&paper).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "提交失败"})
		return
	}

	for _, item := range req.Answer {
		if item.Qid <= 0 || strings.TrimSpace(item.Answer) == "" {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{Code: 500, Msg: "答案格式错误"})
			return
		}
		var question models.GreenQuestion
		if err := tx.First(&question, item.Qid).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{Code: 500, Msg: "题目不存在"})
			return
		}
		answer := models.GreenPaperAnswer{
			PaperId:     paper.ID,
			QuestionId:  question.ID,
			UserAnswer:  strings.TrimSpace(item.Answer),
			RightAnswer: question.Answer,
			IsCorrect:   map[bool]string{true: "1", false: "0"}[strings.EqualFold(strings.TrimSpace(item.Answer), strings.TrimSpace(question.Answer))],
		}
		if err := tx.Create(&answer).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{Code: 500, Msg: "提交失败"})
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "提交失败"})
		return
	}

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "请求成功"})
}

func DataList1(c *gin.Context) {
	listGreenData(c, "list_1")
}

func DataList2(c *gin.Context) {
	listGreenData(c, "list_2")
}

func DataList3(c *gin.Context) {
	listGreenData(c, "list_3")
}

func DataList4(c *gin.Context) {
	listGreenData(c, "list_4")
}

func listGreenData(c *gin.Context, listKey string) {
	var rows []models.GreenDataSeries
	if err := config.DB.Where("list_key = ?", listKey).Order("sort asc, id asc").Find(&rows).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "查询失败"})
		return
	}

	resp := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, gin.H{
			"name": row.Name,
			"data": parseGreenDataValue(row.Data),
		})
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "请求成功", Data: resp})
}

func parseGreenDataValue(raw string) interface{} {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	var parsed interface{}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
		return parsed
	}
	return raw
}

func GreenDataCardCreate(c *gin.Context) {
	var req struct {
		Icon  string `json:"icon"`
		Title string `json:"title"`
		Num   string `json:"num"`
		Unit  string `json:"unit"`
		Trend string `json:"trend"`
		Sort  int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	card := models.GreenDataCard{
		Icon:  req.Icon,
		Title: req.Title,
		Num:   req.Num,
		Unit:  req.Unit,
		Trend: req.Trend,
		Sort:  req.Sort,
	}
	if err := config.DB.Create(&card).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "创建失败"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "创建成功"})
}

func GreenDataCardUpdate(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Icon  string `json:"icon"`
		Title string `json:"title"`
		Num   string `json:"num"`
		Unit  string `json:"unit"`
		Trend string `json:"trend"`
		Sort  int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if err := config.DB.Model(&models.GreenDataCard{}).Where("id = ?", id).Updates(map[string]interface{}{
		"icon":  req.Icon,
		"title": req.Title,
		"num":   req.Num,
		"unit":  req.Unit,
		"trend": req.Trend,
		"sort":  req.Sort,
	}).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "更新失败"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "更新成功"})
}

func GreenDataCardDelete(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.GreenDataCard{}, id).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "删除失败"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func GreenQuestionList(c *gin.Context) {
	var list []models.GreenQuestion
	p := int64(1)
	s := int64(10)
	if v, err := strconv.ParseInt(c.DefaultQuery("pageNum", "1"), 10, 64); err == nil && v > 0 {
		p = v
	}
	if v, err := strconv.ParseInt(c.DefaultQuery("pageSize", "10"), 10, 64); err == nil && v > 0 {
		s = v
	}
	offset := (p - 1) * s
	query := config.DB.Model(&models.GreenQuestion{})
	var total int64
	query.Count(&total)
	if err := query.Offset(int(offset)).Limit(int(s)).Order("id desc").Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "查询失败"})
		return
	}
	respondList(c, "请求成功", list, total)
}

func GreenQuestionCreate(c *gin.Context) {
	var req models.GreenQuestion
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "创建失败"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "创建成功"})
}

func GreenQuestionUpdate(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		QuestionType string `json:"questionType"`
		Level        string `json:"level"`
		Question     string `json:"question"`
		OptionA      string `json:"optionA"`
		OptionB      string `json:"optionB"`
		OptionC      string `json:"optionC"`
		OptionD      string `json:"optionD"`
		OptionE      string `json:"optionE"`
		OptionF      string `json:"optionF"`
		Answer       string `json:"answer"`
		Score        int    `json:"score"`
		Status       string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if err := config.DB.Model(&models.GreenQuestion{}).Where("id = ?", id).Updates(map[string]interface{}{
		"question_type": req.QuestionType,
		"level":         req.Level,
		"question":      req.Question,
		"option_a":      req.OptionA,
		"option_b":      req.OptionB,
		"option_c":      req.OptionC,
		"option_d":      req.OptionD,
		"option_e":      req.OptionE,
		"option_f":      req.OptionF,
		"answer":        req.Answer,
		"score":         req.Score,
		"status":        req.Status,
	}).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "更新失败"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "更新成功"})
}

func GreenQuestionDelete(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.GreenQuestion{}, id).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "删除失败"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func GreenDataSeriesByKey(c *gin.Context) {
	listKey := c.Param("listKey")
	if listKey == "" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	var record models.GreenDataSeries
	if err := config.DB.Where("list_key = ?", listKey).First(&record).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 404, Msg: "数据不存在"})
		return
	}

	var parsed []map[string]interface{}
	if err := json.Unmarshal([]byte(record.Data), &parsed); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "数据解析失败"})
		return
	}

	type DataItem struct {
		Name string      `json:"name"`
		Data interface{} `json:"data"`
	}
	result := make([]DataItem, 0, len(parsed))
	for _, item := range parsed {
		result = append(result, DataItem{
			Name: item["name"].(string),
			Data: item["data"],
		})
	}

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "请求成功", Data: result})
}

func GreenDataSeriesList(c *gin.Context) {
	var list []models.GreenDataSeries
	if err := config.DB.Order("CAST(SUBSTR(list_key, 6) AS INTEGER) asc, id asc").Find(&list).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "查询失败"})
		return
	}
	respondList(c, "请求成功", list, int64(len(list)))
}

func GreenDataSeriesCreate(c *gin.Context) {
	var req struct {
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if req.Data == "" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "数据不能为空"})
		return
	}

	for i := 0; i < 5; i++ {
		var maxNum int
		rows, err := config.DB.Raw("SELECT list_key FROM green_data_series WHERE list_key LIKE 'list_%' ORDER BY CAST(SUBSTR(list_key, 5) AS INTEGER) DESC LIMIT 1").Rows()
		if err != nil {
			c.JSON(http.StatusOK, Response{Code: 500, Msg: "创建失败"})
			return
		}
		if rows.Next() {
			var key string
			if scanErr := rows.Scan(&key); scanErr == nil {
				_, _ = fmt.Sscanf(key, "list_%d", &maxNum)
			}
		}
		rows.Close()

		newListKey := fmt.Sprintf("list_%d", maxNum+1)
		record := models.GreenDataSeries{
			ListKey: newListKey,
			Name:    newListKey,
			Data:    req.Data,
			Sort:    maxNum + 1,
		}
		if err := config.DB.Create(&record).Error; err != nil {
			if isUniqueConstraintError(err) {
				continue
			}
			c.JSON(http.StatusOK, Response{Code: 500, Msg: "创建失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": gin.H{"listKey": newListKey}})
		return
	}

	c.JSON(http.StatusOK, Response{Code: 500, Msg: "创建失败，请重试"})
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint") || strings.Contains(msg, "duplic")
}

func GreenDataSeriesUpdate(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if req.Data == "" {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "数据不能为空"})
		return
	}
	if err := config.DB.Model(&models.GreenDataSeries{}).Where("id = ?", id).Update("data", req.Data).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "更新失败"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "更新成功"})
}

func GreenDataSeriesDelete(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.GreenDataSeries{}, id).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "删除失败"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

func generateToken(userId int, userName, nickName, phone string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"userId":   userId,
		"userName": userName,
		"nickName": nickName,
		"phone":    phone,
		"iat":      now.Unix(),
		"nbf":      now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("digital-community-secret-key-2024"))
}
