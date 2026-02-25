package router

import (
	"digital-community/internal/handlers"
	"digital-community/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	r.Static("/profile", "./profile")
	r.POST("/logout", handlers.Logout)

	prodApi := r.Group("/prod-api/api")
	{
		prodApi.POST("/phone/login", handlers.PhoneLogin)
		prodApi.POST("/login", handlers.Login)
		prodApi.GET("/SMSCode", handlers.SMSCode)
		prodApi.POST("/register", handlers.Register)

		prodApi.POST("/logout", handlers.Logout)

		prodApi.GET("/rotation/list", handlers.RotationList)

		prodApi.GET("/press/category/list", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressCategoryList)
		prodApi.GET("/press/newsList", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressNewsList)
		prodApi.GET("/press/category/newsList", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressCategoryNewsList)
		prodApi.GET("/press/news/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressNewsDetail)
		prodApi.PUT("/press/like/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressLike)

		prodApi.POST("/comment/pressComment", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressComment)
		prodApi.GET("/comment/comment/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.CommentList)
		prodApi.PUT("/comment/like/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.CommentLike)

		prodApi.POST("/common/upload", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.Upload)

		prodApi.GET("/notice/list", handlers.NoticeList)
		prodApi.GET("/notice/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.NoticeDetail)
		prodApi.PUT("/readNotice/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.ReadNotice)

		prodApi.GET("/friendly_neighborhood/list", handlers.FriendlyNeighborList)
		prodApi.POST("/friendly_neighborhood/add/comment", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.FriendlyNeighborAddComment)
		prodApi.GET("/friendly_neighborhood/:id", handlers.FriendlyNeighborDetail)

		prodApi.GET("/activity/topList", handlers.ActivityTopList)
		prodApi.GET("/activity/List", handlers.ActivityList)
		prodApi.POST("/activity/search", handlers.ActivitySearch)
		prodApi.GET("/activity/:id", handlers.ActivityDetail)
		prodApi.GET("/activity/category/list/:id", handlers.ActivityCategoryList)

		prodApi.POST("/registration", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.Registration)
		prodApi.PUT("/checkin/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.Checkin)
		prodApi.PUT("/registration/comment/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.RegistrationComment)

		prodApi.GET("/user/getUserInfo", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GetUserInfo)
		prodApi.PUT("/user/updateUserInfo", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.UpdateUserInfo)
		prodApi.PUT("/user/resetPwd", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.ResetPwd)
	}

	return r
}
