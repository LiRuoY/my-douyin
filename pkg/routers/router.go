package routers

import (
	"douyin/config"
	"douyin/handlers"
	"douyin/pkg/myjwt"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine) {

	//r.Static(constants.StaticHTTPFolder, constants.StaticLocalFolderPath)
	r.Static(config.CStatic.HttpFolder, config.CStatic.LocalStorePath)

	douyin1 := r.Group("/douyin")
	noAuthority := douyin1.Group("/")
	{
		noAuthority.POST("/user/login/", handlers.UserLogin)
		noAuthority.POST("/user/register/", handlers.UserRegister)

	}

	authority := douyin1.Group("/")
	authority.Use(myjwt.JWTMiddleWareImpl)
	{
		{
			authority.GET("/user/", handlers.UserInfo)
		}
		{
			authority.POST("/publish/action/", handlers.PublicVideo)
			authority.GET("/publish/list/", handlers.VideoList)
		}
		{
			authority.GET("/feed", handlers.FeedFlow)
		}
		{
			authority.POST("/favorite/action/", handlers.FavorVideo)
			authority.GET("/favorite/list/", handlers.FavorVideoList)
		}
		{
			authority.POST("/comment/action/", handlers.CommentVideo)
			authority.GET("/comment/list/", handlers.CommentListOfVideo)
		}
		{
			authority.POST("/relation/action/", handlers.FollowUser)
			authority.GET("/relation/follow/list/", handlers.FollowList)
			authority.GET("/relation/follower/list/", handlers.FollowerList)
		}
	}

}
