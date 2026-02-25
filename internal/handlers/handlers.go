package handlers

import (
	"digital-community/internal/config"
	"digital-community/internal/models"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	Phonenumber  string  `json:"phonenumber"`
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
		Phonenumber:  user.Phone,
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
		"subTitle":    "",
		"content":     news.Content,
		"publishDate": publishDate,
		"commentNum":  0,
		"likeNum":     0,
		"readNum":     news.ViewCount,
		"top":         "N",
		"hot":         "N",
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
		"isTop":         "1",
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
		Phone   string `json:"phone" binding:"required"`
		SMSCode string `json:"SMSCode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}
	if !verifySMSCode(req.Phone, req.SMSCode) {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "验证码错误或已过期"})
		return
	}

	var user models.User
	result := config.DB.Where("phone = ?", req.Phone).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		user = models.User{
			Phone:    req.Phone,
			NickName: "用户" + req.Phone[len(req.Phone)-4:],
			Status:   "0",
			DelFlag:  "0",
		}
		config.DB.Create(&user)
		config.DB.Where("phone = ?", req.Phone).First(&user)
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
		UserName string `json:"userName" binding:"required"`
		PassWord string `json:"passWord" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	var user models.User
	result := config.DB.Where("user_name = ? AND pass_word = ?", req.UserName, req.PassWord).First(&user)
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
		PassWord     string `json:"passWord" binding:"required"`
		Phonenumber  string `json:"phonenumber"`
		Phone        string `json:"phone"`
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

	var count int64
	config.DB.Model(&models.User{}).Where("user_name = ?", req.UserName).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "用户名已存在"})
		return
	}

	user := models.User{
		UserName:     req.UserName,
		PassWord:     req.PassWord,
		NickName:     req.NickName,
		Phone:        req.Phone,
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
	if req.Phonenumber != "" {
		user.Phone = req.Phonenumber
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
		NickName     string `json:"nickName"`
		Phonenumber  string `json:"phonenumber"`
		Phone        string `json:"phone"`
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

	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "用户不存在"})
		return
	}

	config.DB.Model(&user).Updates(map[string]interface{}{
		"nick_name":    req.NickName,
		"phone":        req.Phone,
		"sex":          req.Sex,
		"avatar":       req.Avatar,
		"email":        req.Email,
		"id_card":      req.IDCard,
		"address":      req.Address,
		"introduction": req.Introduction,
	})
	if req.Phonenumber != "" {
		config.DB.Model(&user).Update("phone", req.Phonenumber)
	}

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func ResetPwd(c *gin.Context) {
	userId := c.GetInt("userId")
	var req struct {
		PassWord string `json:"passWord" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	config.DB.Model(&models.User{}).Where("id = ?", userId).Update("pass_word", req.PassWord)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "退出成功"})
}

func RotationList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)
	rtype := c.Query("type")

	var rotations []models.Rotation
	query := config.DB.Model(&models.Rotation{})
	if rtype != "" {
		query = query.Where("type = ?", rtype)
	}

	var total int64
	query.Count(&total)

	query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&rotations)

	items := make([]gin.H, 0, len(rotations))
	for _, r := range rotations {
		items = append(items, gin.H{
			"id":     r.ID,
			"advImg": r.PicPath,
			"type":   strconv.Itoa(r.Type),
		})
	}
	respondList(c, "请求成功", items, total)
}

func PressCategoryList(c *gin.Context) {
	var categories []models.PressCategory
	config.DB.Find(&categories)
	items := make([]gin.H, 0, len(categories))
	for _, v := range categories {
		items = append(items, gin.H{"id": v.ID, "name": v.Name, "sort": v.Sort})
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "请求成功", Data: items})
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

	var newsList []models.PressNews
	query := config.DB.Model(&models.PressNews{})
	if id != "" {
		query = query.Where("category_id = ?", id)
	}

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

func PressLike(c *gin.Context) {
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func PressComment(c *gin.Context) {
	var req struct {
		Sid     int    `json:"sid" binding:"required"`
		Content string `json:"content" binding:"required"`
		Type    string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	userId := c.GetInt("userId")
	nickName := c.GetString("nickName")

	comment := models.Comment{
		Sid:        req.Sid,
		Type:       req.Type,
		Content:    req.Content,
		UserId:     userId,
		NickName:   nickName,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	config.DB.Create(&comment)

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
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "上传文件失败"})
		return
	}

	filename := "/profile/upload/" + time.Now().Format("20060102150405") + "_" + file.Filename
	c.SaveUploadedFile(file, "."+filename)

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "操作成功",
		"fileName": file.Filename,
		"url":      filename,
	})
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

func ReadNotice(c *gin.Context) {
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
		NeighborId int    `json:"neighborId" binding:"required"`
		Content    string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	userId := c.GetInt("userId")
	nickName := c.GetString("nickName")

	comment := models.FNComment{
		NeighborId: req.NeighborId,
		UserId:     userId,
		NickName:   nickName,
		Content:    req.Content,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	config.DB.Create(&comment)

	config.DB.Model(&models.FriendlyNeighbor{}).Where("id = ?", req.NeighborId).UpdateColumn("comment_num", gorm.Expr("comment_num + ?", 1))

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
			"pulishTime":     v.CreateTime,
			"neighborhoodId": v.NeighborId,
		})
	}
	item := buildNeighborItem(neighbor)
	item["userComment"] = commentItems
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "请求成功", Data: item})
}

func ActivityTopList(c *gin.Context) {
	pageNum, pageSize := parsePaging(c)

	var activities []models.Activity
	var total int64
	config.DB.Model(&models.Activity{}).Count(&total)

	config.DB.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("start_date DESC").Find(&activities)
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
		Words string `json:"words"`
	}
	_ = c.ShouldBindJSON(&req)

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

func Registration(c *gin.Context) {
	var req struct {
		ActivityId int    `json:"activityId" binding:"required"`
		Phone      string `json:"phone" binding:"required"`
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

	registration := models.Registration{
		UserId:     userId,
		UserName:   userName,
		NickName:   nickName,
		Phone:      req.Phone,
		ActivityId: req.ActivityId,
		Status:     "0",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	config.DB.Create(&registration)

	config.DB.Model(&models.Activity{}).Where("id = ?", req.ActivityId).UpdateColumn("current_count", gorm.Expr("current_count + ?", 1))

	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func Checkin(c *gin.Context) {
	id := c.Param("id")
	config.DB.Model(&models.Registration{}).Where("id = ?", id).Update("checkin_status", "1")
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func RegistrationComment(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Comment string `json:"comment" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Response{Code: 500, Msg: "参数错误"})
		return
	}

	config.DB.Model(&models.Registration{}).Where("id = ?", id).Update("comment", req.Comment)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "操作成功"})
}

func generateToken(userId int, userName, nickName, phone string) (string, error) {
	claims := jwt.MapClaims{
		"userId":   userId,
		"userName": userName,
		"nickName": nickName,
		"phone":    phone,
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("digital-community-secret-key-2024"))
}
