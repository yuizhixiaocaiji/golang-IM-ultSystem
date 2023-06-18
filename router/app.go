package router

import (
	"ginchat/docs"
	"ginchat/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	docs.SwaggerInfo.BasePath = ""
	r.GET("api/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//静态资源
	r.Static("api/assets", "assets/")

	r.POST("api/searchFriends", service.SearchFriends)

	r.GET("api/index", service.GetIndex)
	r.GET("api/user/getUserList", service.GetUserList)
	r.GET("api/user/createUser", service.CreateUser)
	r.GET("api/user/deleteUser", service.DeleteUser)
	r.GET("api/chat", service.Chat)
	r.POST("api/user/updateUser", service.UpdateUser)
	r.POST("api/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)

	//发送消息
	r.GET("api/user/sendMsg", service.SendMsg)
	//发送消息
	r.GET("api/user/sendUserMsg", service.SendUserMsg)
	//上传文件
	r.POST("api/attach/upload", service.Upload)
	//添加好友
	r.POST("api/contact/addFriend", service.AddFriend)
	//创建群
	r.POST("api/contact/community", service.CreateCommunity)
	return r
}
