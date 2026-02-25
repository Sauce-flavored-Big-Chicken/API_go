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
		prodApi.GET("/smsCode", handlers.SMSCode)
		prodApi.POST("/register", handlers.Register)

		prodApi.GET("/rotation/list", handlers.RotationList)
		prodApi.POST("/rotation", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.RotationCreate)
		prodApi.PUT("/rotation/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.RotationUpdate)
		prodApi.DELETE("/rotation/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.RotationDelete)

		prodApi.GET("/press/category/list", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressCategoryList)
		prodApi.POST("/press/category", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressCategoryCreate)
		prodApi.PUT("/press/category/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressCategoryUpdate)
		prodApi.DELETE("/press/category/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressCategoryDelete)

		prodApi.GET("/press/newsList", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressNewsList)
		prodApi.POST("/press/news", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressNewsCreate)
		prodApi.PUT("/press/news/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressNewsUpdate)
		prodApi.DELETE("/press/news/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressNewsDelete)
		prodApi.GET("/press/category/newsList", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressCategoryNewsList)
		prodApi.GET("/press/news/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressNewsDetail)
		prodApi.PUT("/press/like/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressLike)

		prodApi.POST("/comment/pressComment", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.PressComment)
		prodApi.GET("/comment/comment/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.CommentList)
		prodApi.PUT("/comment/like/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.CommentLike)

		prodApi.POST("/common/upload", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.Upload)
		prodApi.GET("/common/datacard", handlers.CommonDataCard)
		prodApi.POST("/common/datacard", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GreenDataCardCreate)
		prodApi.PUT("/common/datacard/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GreenDataCardUpdate)
		prodApi.DELETE("/common/datacard/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GreenDataCardDelete)
		prodApi.GET("/common/images", handlers.ImageList)
		prodApi.DELETE("/common/images", handlers.ImageDelete)
		prodApi.GET("/common/files", handlers.FileList)
		prodApi.DELETE("/common/files", handlers.FileDelete)

		prodApi.GET("/question/questionList/:id/:level", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.QuestionQuestionList)
		prodApi.POST("/question/savePaper", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.QuestionSavePaper)
		prodApi.GET("/question/list", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GreenQuestionList)
		prodApi.POST("/question", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GreenQuestionCreate)
		prodApi.PUT("/question/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GreenQuestionUpdate)
		prodApi.DELETE("/question/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GreenQuestionDelete)

		prodApi.GET("/data/list", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GreenDataSeriesList)
		prodApi.GET("/data/:listKey", handlers.GreenDataSeriesByKey)
		prodApi.POST("/data/list", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GreenDataSeriesCreate)
		prodApi.PUT("/data/list/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GreenDataSeriesUpdate)
		prodApi.DELETE("/data/list/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GreenDataSeriesDelete)

		prodApi.GET("/notice/list", handlers.NoticeList)
		prodApi.POST("/notice", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.NoticeCreate)
		prodApi.PUT("/notice/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.NoticeUpdate)
		prodApi.DELETE("/notice/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.NoticeDelete)
		prodApi.GET("/notice/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.NoticeDetail)
		prodApi.PUT("/readNotice/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.ReadNotice)

		prodApi.GET("/friendly_neighborhood/list", handlers.FriendlyNeighborList)
		prodApi.POST("/friendly_neighborhood", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.FriendlyNeighborCreate)
		prodApi.PUT("/friendly_neighborhood/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.FriendlyNeighborUpdate)
		prodApi.DELETE("/friendly_neighborhood/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.FriendlyNeighborDelete)
		prodApi.POST("/friendly_neighborhood/add/comment", handlers.FriendlyNeighborAddComment)
		prodApi.GET("/friendly_neighborhood/:id", handlers.FriendlyNeighborDetail)

		prodApi.GET("/activity/topList", handlers.ActivityTopList)
		prodApi.GET("/activity/list", handlers.ActivityList)
		prodApi.POST("/activity/search", handlers.ActivitySearch)
		prodApi.GET("/activity/category/list/:id", handlers.ActivityCategoryList)
		prodApi.GET("/activity/:id", handlers.ActivityDetail)
		prodApi.POST("/activity", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.ActivityCreate)
		prodApi.PUT("/activity/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.ActivityUpdate)
		prodApi.DELETE("/activity/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.ActivityDelete)

		prodApi.GET("/registration/list", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.RegistrationList)
		prodApi.POST("/registration", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.Registration)
		prodApi.PUT("/checkin/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.Checkin)
		prodApi.PUT("/registration/comment/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.RegistrationComment)

		prodApi.GET("/user/list", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.UserList)
		prodApi.GET("/user/getUserInfo", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.GetUserInfo)
		prodApi.POST("/user", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.UserCreate)
		prodApi.PUT("/user/updateUserInfo", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.UpdateUserInfo)
		prodApi.PUT("/user/resetPwd", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.ResetPwd)
		prodApi.PUT("/user/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.UserUpdate)
		prodApi.DELETE("/user/:id", middleware.AuthMiddleware("digital-community-secret-key-2024"), handlers.UserDelete)
	}

	return r
}
