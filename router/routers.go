package router

import (
	"Umeet/controller"
	"Umeet/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router() {
	// gin.DisableConsoleColor()
	// f, _ := os.Create("umeet.log")
	// gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	// r.Use(middleware.IPLimiter)

	common := r.RouterGroup
	{
		common.POST("/Login", controller.Login)          //登录
		common.POST("/Register", controller.Register)    //注册
		common.GET("/VerifyUser", controller.VerifyUser) //验证
		common.StaticFS("/static", http.Dir("static"))
		common.GET("/UpdateToken", controller.UpdateToken)
	}

	usermanager := r.Group("user").Use(middleware.Authorization)
	{
		usermanager.POST("/UpdateAvatar", controller.UpdateAvatar) //更新用户头像
		usermanager.PUT("EditProfile", controller.UpdatesUser)     //更改用户信息
	}

	postmanager := r.Group("post").Use(middleware.Authorization)
	{
		postmanager.GET("/GetCustomerPosts", controller.GetUserPost) //获取用户发布的帖子
		postmanager.POST("/AddPicture", controller.AddPicture)       //添加图片
		postmanager.GET("/GetHomePosts", controller.GetHomePost)     //获取主页帖子
		postmanager.POST("/AddPost", controller.AddPost)             //添加帖子
		postmanager.GET("/LikeorNot", controller.LikePostByRedis)    //点赞帖子
		postmanager.DELETE("/deletePost", controller.DeletePost)     //删除帖子
		postmanager.GET("/GetPostImages", controller.GetPostImages)  //获取帖子的图片
		postmanager.POST("/AddTaxiPost", controller.AddCarPost)      //添加拼车帖子
		postmanager.GET("/GetTaxiPosts", controller.AllCarPosts)     //获得所有拼车帖子
		postmanager.PUT("/joinTaxi", controller.JoinCarSharing)      //加入拼车
	}

	commentmanager := r.Group("comment").Use(middleware.Authorization)
	{
		commentmanager.POST("/AddRootComment", controller.AddRootComment)      //添加根评论
		commentmanager.POST("/AddChildComment", controller.AddChildComment)    //添加子评论
		commentmanager.DELETE("/DeleteComment", controller.DeleteComment)      //删除评论
		commentmanager.GET("/getComments", controller.GetRootComments)         //获取根评论
		commentmanager.GET("/getChildrenComments", controller.GetChildComment) //获取子评论
		commentmanager.GET("/LikeorNot", controller.LikeComment)               //点赞评论
	}

	administrator := r.Group("admin")
	{
		administrator.GET("/posts", controller.AllPosts)
		administrator.GET("/users", controller.AllUsers)
		administrator.DELETE("/post", controller.DeletePost)
		administrator.GET("/comments", controller.AllComment)
		administrator.DELETE("/comment", controller.DeleteComment)
		administrator.POST("/banuser", controller.DeleteUser)
	}
	r.POST("/imageChat", controller.ImageChat)
	r.GET("/room", controller.SelectRoom)
	r.Run("0.0.0.0:8080")
	// r.Run(config.IP["addr"])

}

// clientManager := websocket.NewManager()
// websocket.Rooms["1"] = clientManager
// go clientManager.Run()
// r.GET("/room", func(ctx *gin.Context) {
// 	name := ctx.Query("name")
// 	uid := ctx.Query("uid")
// 	websocket.ServeWs(clientManager, ctx.Writer, ctx.Request, name, uid)
// })

// wsmanager := r.Group("ws")
// {
// 	wsmanager.GET("/CreateRoom", controller.CreateRoom)
// }
