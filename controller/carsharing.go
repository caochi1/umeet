package controller

import (
	"Umeet/models"
	"Umeet/websocket"

	"github.com/gin-gonic/gin"
)

type JsResponse struct {
	Code int `json:"code"`
	Data any `json:"data"`
	Msg  any `json:"msg"`
}

func Response(ctx *gin.Context, code int, data any, msg any) {
	ctx.JSON(code, JsResponse{code, data, msg})
}

// 添加拼车帖子+创建聊天室
func AddCarPost(ctx *gin.Context) {
	var cs models.CarSharing
	if ctx.ShouldBindJSON(&cs) != nil {
		Response(ctx, 404, nil, "error")
		return
	}
	if models.AddCarPost(&cs) != nil {
		Response(ctx, 404, nil, "error")
		return
	}
	Response(ctx, 200, nil, nil)

}

// 加入拼车
func JoinCarSharing(ctx *gin.Context) {
	uid, pid, nickname := ctx.Query("uid"), ctx.Query("pid"), ctx.Query("nickname")
	room := websocket.Rooms.Get(pid)
	if room == nil {
		Response(ctx, 404, nil, "房间不存在")
		return
	}
	if _, ok := room.Get(uid); !ok {
		room.Set(uid, nickname)
	} else {
		room.Delete(uid)
	}
}

// 所有拼车帖子
func AllCarPosts(ctx *gin.Context) {
	result, _ := models.AllCarPosts()
	Response(ctx, 200, &result, nil)
}

// 进入聊天室
func SelectRoom(ctx *gin.Context) {
	uid, rid, name := ctx.Query("uid"), ctx.Query("rid"), ctx.Query("username")
	room := websocket.Rooms.Get(rid)
	if room == nil {
		Response(ctx, 404, nil, "房间不存在")
		return
	}
	websocket.ServeWs(room, ctx.Writer, ctx.Request, name, rid, uid)

}

// 	clientManager := websocket.NewManager()
// 	go clientManager.Run()
// 	R.GET("/room", func(ctx *gin.Context) {
// 		name := ctx.Query("name")
// 		websocket.ServeWs(clientManager, ctx.Writer, ctx.Request, name)
// 	})
// }

// clientManager2 := websocket.NewManager()
// go clientManager2.Run()
// r.GET("/room2", func(ctx *gin.Context) {
// 	name := ctx.Query("name")
// 	websocket.ServeWs(clientManager2, ctx.Writer, ctx.Request, name)
// })
