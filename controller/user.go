package controller

import (
	"Umeet/models"
	"Umeet/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 验证用户是否存在
func VerifyUser(ctx *gin.Context) {
	username := ctx.Query("username")
	if models.VerifyUser(username) == 1 {
		Response(ctx, 404, nil, "账户已存在")
		return
	}
	Response(ctx, 200, nil, "验证成功")
}

// 登录
func Login(ctx *gin.Context) {
	var user models.Userinfo
	if err := ctx.ShouldBindJSON(&user); err != nil {
		Response(ctx, 404, nil, "错误")
		return
	}
	u, err := models.CheckUser(user.UserName)
	if err != nil {
		Response(ctx, 404, nil, "用户名或密码错误")
		return
	}
	if u.Passwd != utils.EncryMd5(user.Passwd) {
		Response(ctx, 404, nil, "用户名或密码错误")
		return
	}
	token, _ := utils.GenerateToken(u.ID, u.NickName)
	Response(ctx, 200, token, u)
}

// 注册
func Register(ctx *gin.Context) {
	var user models.Userinfo
	if err := ctx.ShouldBindJSON(&user); err != nil {
		Response(ctx, 404, nil, "错误")
		return
	}
	user.Passwd = utils.EncryMd5(user.Passwd)
	err := models.CreateUser(&user)
	if err != nil {
		Response(ctx, 404, nil, "未知错误")
	} else {
		token, _ := utils.GenerateToken(user.ID, user.NickName)
		Response(ctx, 200, token, user)
	}
}

// 获取用户发布的帖子
func GetUserPost(ctx *gin.Context) {
	uid := utils.S2uint(ctx.Query("uid"))
	data, err := models.GetPostByUserID(uid)
	if err != nil {
		Response(ctx, 404, nil, "错误")
		return
	}
	Response(ctx, 200, data, "success")
}

// 更新token
func UpdateToken(ctx *gin.Context) {
	uid, username := utils.S2uint(ctx.Query("uid")), ctx.Query("nickname")
	token, err := utils.GenerateToken(uid, username)
	if err != nil {
		Response(ctx, 404, err, "failed")
	}
	Response(ctx, 200, token, "new token")
}

// 更改用户信息
func UpdatesUser(ctx *gin.Context) {
	var user models.Userinfo
	if err := ctx.ShouldBindJSON(&user); err != nil {
		Response(ctx, 404, nil, "错误")
		return
	}
	models.UpdatesUser(user.ID, user)
	Response(ctx, 200, nil, nil)
}

func AllUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("per_page"))
	users, err := models.AllUsers(page, pageSize)
	Response(ctx, 200, users, err)
}

func DeleteUser(ctx *gin.Context) {
	uid := utils.S2uint(ctx.Query("id"))
	models.RmUser(uid)
	Response(ctx, 200, "success", "success")
}

// 获取点赞的帖子
// func GetFavorPosts(ctx *gin.Context) {

// 	data, err := models.GetFavorPostsByUserID(uid.Uid)
// 	if err != nil {
// 		Response(ctx, 404, nil, "错误")
// 		return
// 	}
// 	Response(ctx, 200, data, "success")
// }
