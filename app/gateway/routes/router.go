package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/CocaineCong/grpc-todolist/app/gateway/internal/http"
	"github.com/CocaineCong/grpc-todolist/app/gateway/middleware"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors(), middleware.ErrorMiddleware())
	store := cookie.NewStore([]byte("something-very-secret"))
	ginRouter.Use(sessions.Sessions("mysession", store))
	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})
		// 用户服务
		v1.POST("/user/register", http.UserRegister)
		v1.POST("/user/login", http.UserLogin)

		// 需要登录保护
		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			// 任务模块
			authed.GET("task", http.GetTaskList)
			authed.POST("task", http.CreateTask)
			authed.PUT("task", http.UpdateTask)
			authed.DELETE("task", http.DeleteTask)
			authed.POST("upload",http.UploadFile)
		}
		// 添加抽奖路由
		authed1 := v1.Group("/lottery")
		authed1.Use(middleware.JWT())
		{
			// 任务模块
			authed1.POST("initaward", http.InitAward)   // 初始化奖品
			authed1.POST("draw", http.Draw)   // 初始化奖品
			authed1.GET("listAwardInfo", http.ListAwardInfo)   // 初始化奖品
			authed1.GET("toMysql", http.ToMysql)   // 初始化奖品
			//authed.POST("task", http.CreateTask)
			//authed.PUT("task", http.UpdateTask)
			//authed.DELETE("task", http.DeleteTask)
		}
	}
	return ginRouter
}
