package config

import (
	"digital-community/internal/models"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
	var err error

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Rotation{},
		&models.PressCategory{},
		&models.PressNews{},
		&models.PressLikeRecord{},
		&models.Notice{},
		&models.FriendlyNeighbor{},
		&models.FNComment{},
		&models.Activity{},
		&models.ActivityCategory{},
		&models.Registration{},
		&models.Comment{},
		&models.CommentLikeRecord{},
		&models.GreenDataCard{},
		&models.GreenQuestion{},
		&models.GreenPaper{},
		&models.GreenPaperAnswer{},
		&models.GreenDataSeries{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if err := seedDefaultUser(); err != nil {
		return fmt.Errorf("failed to seed default user: %w", err)
	}

	if err := seedBusinessData(); err != nil {
		return fmt.Errorf("failed to seed business data: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func seedDefaultUser() error {
	var count int64
	if err := DB.Model(&models.User{}).Where("user_name = ?", "test01").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	user := models.User{
		UserName:     "test01",
		NickName:     "测试用户",
		PassWord:     "123456",
		Phone:        "13800000000",
		Sex:          "0",
		Email:        "test01@example.com",
		Status:       "0",
		DelFlag:      "0",
		Address:      "默认地址",
		Introduction: "系统默认用户",
		Balance:      0,
		Score:        0,
		LoginDate:    now,
	}
	return DB.Create(&user).Error
}

func seedBusinessData() error {
	images, err := prepareSeedImages()
	if err != nil {
		return err
	}

	if err := seedUserAvatars(images); err != nil {
		return err
	}
	if err := seedRotations(images); err != nil {
		return err
	}
	if err := seedPressCategories(); err != nil {
		return err
	}
	if err := seedPressNews(images); err != nil {
		return err
	}
	if err := seedNotices(); err != nil {
		return err
	}
	if err := seedFriendlyNeighbors(images); err != nil {
		return err
	}
	if err := seedActivityCategories(); err != nil {
		return err
	}
	if err := seedActivities(images); err != nil {
		return err
	}

	return nil
}

func prepareSeedImages() ([]string, error) {
	srcDir := "./alcy_fj_100"
	dstDir := "./profile/upload/image/seed"

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return nil, err
	}

	if entries, err := os.ReadDir(srcDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			ext := strings.ToLower(filepath.Ext(name))
			switch ext {
			case ".jpg", ".jpeg", ".png", ".webp", ".gif":
			default:
				continue
			}

			src := filepath.Join(srcDir, name)
			dst := filepath.Join(dstDir, name)
			if _, statErr := os.Stat(dst); statErr == nil {
				continue
			}
			if err := copyFile(src, dst); err != nil {
				return nil, err
			}
		}
	}

	dstEntries, err := os.ReadDir(dstDir)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, entry := range dstEntries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp", ".gif":
			files = append(files, "/profile/upload/image/seed/"+name)
		}
	}

	sort.Strings(files)
	return files, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func pickImage(images []string, idx int) string {
	if len(images) == 0 {
		return ""
	}
	if idx < 0 {
		idx = -idx
	}
	return images[idx%len(images)]
}

func seedUserAvatars(images []string) error {
	if len(images) == 0 {
		return nil
	}

	var users []models.User
	if err := DB.Find(&users).Error; err != nil {
		return err
	}

	for i := range users {
		if users[i].Avatar != "" {
			continue
		}
		if err := DB.Model(&users[i]).Update("avatar", pickImage(images, i)).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedRotations(images []string) error {
	var count int64
	if err := DB.Model(&models.Rotation{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		items := []models.Rotation{
			{Title: "社区引导一", PicPath: pickImage(images, 0), Type: 1, Status: "0"},
			{Title: "社区引导二", PicPath: pickImage(images, 1), Type: 1, Status: "0"},
			{Title: "社区主页轮播一", PicPath: pickImage(images, 2), Type: 2, Status: "0"},
			{Title: "社区主页轮播二", PicPath: pickImage(images, 3), Type: 2, Status: "0"},
			{Title: "社区主页轮播三", PicPath: pickImage(images, 4), Type: 2, Status: "0"},
		}
		if err := DB.Create(&items).Error; err != nil {
			return err
		}
	}

	var type2Count int64
	if err := DB.Model(&models.Rotation{}).Where("type = ?", 2).Count(&type2Count).Error; err != nil {
		return err
	}
	if type2Count >= 3 {
		return nil
	}

	need := int(3 - type2Count)
	appendItems := make([]models.Rotation, 0, need)
	for i := 0; i < need; i++ {
		appendItems = append(appendItems, models.Rotation{
			Title:   fmt.Sprintf("社区主页轮播补充%d", i+1),
			PicPath: pickImage(images, 10+i),
			Type:    2,
			Status:  "0",
		})
	}

	return DB.Create(&appendItems).Error
}

func seedPressCategories() error {
	var count int64
	if err := DB.Model(&models.PressCategory{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	items := []models.PressCategory{
		{Name: "今日要闻", Sort: 1, Status: "0"},
		{Name: "社区动态", Sort: 2, Status: "0"},
		{Name: "政策通知", Sort: 3, Status: "0"},
	}
	return DB.Create(&items).Error
}

func seedPressNews(images []string) error {
	var count int64
	if err := DB.Model(&models.PressNews{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var categories []models.PressCategory
	if err := DB.Order("sort asc, id asc").Find(&categories).Error; err != nil {
		return err
	}
	if len(categories) == 0 {
		return nil
	}

	now := time.Now()
	items := make([]models.PressNews, 0, 8)
	for i := 0; i < 8; i++ {
		cat := categories[i%len(categories)]
		items = append(items, models.PressNews{
			Title:       fmt.Sprintf("社区资讯第%d期", i+1),
			SubTitle:    "数字社区每日简报",
			Content:     fmt.Sprintf("<p>这是第%d期社区资讯内容，包含社区公告、活动预告和民生服务信息。</p>", i+1),
			ImageUrls:   pickImage(images, 10+i),
			CategoryId:  int(cat.ID),
			Author:      "社区运营中心",
			Source:      "数字社区",
			ViewCount:   10 + i,
			LikeNum:     i % 5,
			CommentNum:  i % 3,
			Type:        strconv.Itoa((i % 3) + 1),
			Top:         []string{"Y", "N"}[i%2],
			Hot:         []string{"Y", "N"}[(i+1)%2],
			Tags:        "社区,民生",
			Status:      "0",
			PublishDate: now.Add(time.Duration(-i) * 24 * time.Hour),
		})
	}

	return DB.Create(&items).Error
}

func seedNotices() error {
	var count int64
	if err := DB.Model(&models.Notice{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now()
	items := []models.Notice{
		{Title: "缴费通知", NoticeStatus: "1", NoticeContent: "请于本月25日前完成物业费缴纳。", PublishDate: now.Add(-48 * time.Hour), CreateBy: "物业管理"},
		{Title: "停水通知", NoticeStatus: "0", NoticeContent: "因管网维护，明日9:00-17:00临时停水。", PublishDate: now.Add(-24 * time.Hour), CreateBy: "社区服务中心"},
		{Title: "消防演练", NoticeStatus: "0", NoticeContent: "本周六上午10点开展消防演练，请居民配合。", PublishDate: now.Add(-6 * time.Hour), CreateBy: "社区网格站"},
	}
	return DB.Create(&items).Error
}

func seedFriendlyNeighbors(images []string) error {
	var count int64
	if err := DB.Model(&models.FriendlyNeighbor{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	items := make([]models.FriendlyNeighbor, 0, 6)
	now := time.Now()
	for i := 0; i < 6; i++ {
		items = append(items, models.FriendlyNeighbor{
			UserId:     1,
			NickName:   fmt.Sprintf("邻里用户%d", i+1),
			UserImgUrl: pickImage(images, 30+i),
			Content:    fmt.Sprintf("这是第%d条友邻动态，欢迎大家互动交流。", i+1),
			ImgUrl:     pickImage(images, 40+i),
			CommentNum: 0,
			LikeNum:    i,
			CreateTime: now.Add(time.Duration(-i) * time.Hour).Format("2006-01-02 15:04:05"),
		})
	}

	return DB.Create(&items).Error
}

func seedActivityCategories() error {
	var count int64
	if err := DB.Model(&models.ActivityCategory{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	items := []models.ActivityCategory{
		{Name: "文化", Status: "0"},
		{Name: "体育", Status: "0"},
		{Name: "公益", Status: "0"},
		{Name: "亲子", Status: "0"},
	}
	return DB.Create(&items).Error
}

func seedActivities(images []string) error {
	var count int64
	if err := DB.Model(&models.Activity{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now()
	items := []models.Activity{
		{
			Title:        "社区植绿行动",
			Content:      "周末一起参与社区绿化，建设美好家园。",
			PicPath:      pickImage(images, 60),
			CategoryId:   3,
			StartDate:    now.Add(7 * 24 * time.Hour),
			EndDate:      now.Add(7*24*time.Hour + 2*time.Hour),
			Address:      "社区中央广场",
			TotalCount:   100,
			CurrentCount: 6,
			IsTop:        "1",
			Status:       "0",
			CreateBy:     "社区服务中心",
			CreateTime:   now.Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
		},
		{
			Title:        "亲子阅读日",
			Content:      "亲子共读，培养阅读习惯。",
			PicPath:      pickImage(images, 61),
			CategoryId:   4,
			StartDate:    now.Add(10 * 24 * time.Hour),
			EndDate:      now.Add(10*24*time.Hour + 3*time.Hour),
			Address:      "社区图书角",
			TotalCount:   40,
			CurrentCount: 12,
			IsTop:        "1",
			Status:       "0",
			CreateBy:     "社区文化站",
			CreateTime:   now.Add(-12 * time.Hour).Format("2006-01-02 15:04:05"),
		},
		{
			Title:        "全民健步走",
			Content:      "倡导健康生活方式，一起健步走。",
			PicPath:      pickImage(images, 62),
			CategoryId:   2,
			StartDate:    now.Add(14 * 24 * time.Hour),
			EndDate:      now.Add(14*24*time.Hour + 2*time.Hour),
			Address:      "滨河步道",
			TotalCount:   200,
			CurrentCount: 35,
			IsTop:        "0",
			Status:       "0",
			CreateBy:     "社区体育组",
			CreateTime:   now.Add(-8 * time.Hour).Format("2006-01-02 15:04:05"),
		},
	}
	return DB.Create(&items).Error
}
