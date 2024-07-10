package controller

import (
	"Umeet/config"
	"Umeet/models"
	"Umeet/utils"
	"Umeet/websocket"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 上传帖子图片
func AddPicture(ctx *gin.Context) {
	spid := ctx.Query("pid")
	pid := utils.S2uint(spid)
	form, err := ctx.MultipartForm()
	if err != nil {
		Response(ctx, 404, nil, "update fail")
		return
	}
	files := form.File["images"]
	dst := utils.StringsBuilder("static/post/", spid, "_")
	for i, file := range files {
		form := strings.Split(file.Filename, ".")
		path := utils.StringsBuilder(dst, strconv.Itoa(i), ".", form[len(form)-1])
		ctx.SaveUploadedFile(file, path)
		models.SavePostPath(utils.StringsBuilder("http://", config.IP["addr"], "/", path), pid)
	}
	Response(ctx, 200, nil, nil)
}

// 更改用户头像
func UpdateAvatar(ctx *gin.Context) {
	uid := ctx.Query("uid")
	file, err := ctx.FormFile("image")
	if err != nil {
		Response(ctx, 404, nil, "update fail")
		return
	}
	id := utils.S2uint(uid)
	user, _ := models.GetUserInfo(id, "avatar")
	oldform := strings.Split(user.Avatar, ".")
	newform := strings.Split(file.Filename, ".")
	index := len(newform) - 1
	dst := utils.StringsBuilder("http://", config.IP["addr"], "/static/user/", uid, ".", newform[index])
	if ctx.SaveUploadedFile(file, utils.StringsBuilder("static/user/", uid, ".", newform[index])) != nil {
		Response(ctx, 404, nil, "update fail")
		return
	}
	if oldform[len(oldform)-1] != newform[index] {
		models.UpdateUser(id, "avatar", dst)
		os.Remove(utils.StringsBuilder("static/user/", uid, ".", oldform[len(oldform)-1]))
	}
	Response(ctx, 200, dst, nil)
}

// 聊天图片
func ImageChat(ctx *gin.Context) {
	uid, name, rid := ctx.Query("uid"), ctx.Query("name"), ctx.Query("rid")
	file, err := ctx.FormFile("image")
	if err != nil {
		Response(ctx, 404, nil, "fail")
		return
	}
	dst := utils.StringsBuilder("static/chat/", strconv.Itoa(int(time.Now().Unix())))
	if ctx.SaveUploadedFile(file, dst) != nil {
		Response(ctx, 404, nil, "fail")
		return
	}
	if cm := websocket.Rooms.Get(rid); cm == nil {
		Response(ctx, 404, nil, "房间不存在")
		return
	} else {
		msg := &websocket.Message{
			Uid:         uid,
			ImgName:     file.Filename,
			Content:     utils.StringsBuilder("http://192.168.2.116:8080/", dst),
			MessageType: 2,
			Size:        file.Size,
			NickName:    name,
		}
		cm.Send(msg)
	}

	// websocket.Rooms.Get(rid).Broadcast <- &websocket.Message{
	// 	Uid:         uid,
	// 	ImgName:     file.Filename,
	// 	Content:     utils.StringsBuilder("http://192.168.2.116:8080/", dst),
	// 	MessageType: 2,
	// 	Size:        file.Size,
	// 	NickName:    name,
	// }

}

// 获取帖子图片
func GetPostImages(ctx *gin.Context) {
	pid := utils.S2uint(ctx.Query("pid"))
	images, err := models.GetPostImages(pid)
	if err != nil {
		Response(ctx, 404, err, nil)
		return
	}
	Response(ctx, 200, images, nil)
}

