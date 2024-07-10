package controller

import (
	"Umeet/cache"
	"Umeet/models"
	"Umeet/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 增加一个帖子
func AddPost(ctx *gin.Context) {
	var p models.Post
	if ctx.ShouldBindJSON(&p) != nil {
		Response(ctx, 404, nil, "error")
		return
	}
	pid, err := models.AddPost(&p)
	if err != nil {
		Response(ctx, 404, nil, err)
		return
	}
	Response(ctx, 200, pid, "AddPost sucess")
}

// 删除一个帖子
func DeletePost(ctx *gin.Context) {
	pid := utils.S2uint(ctx.Query("pid"))
	if err := models.DeletePost(pid); err != nil {
		Response(ctx, 404, nil, err)
		return
	}

	Response(ctx, 200, pid, "DeletePost sucess")
}

// 获取homepage文章
func GetHomePost(ctx *gin.Context) {
	posts := models.GetPostsByRandom(utils.S2uint(ctx.Query("uid")))
	Response(ctx, 200, &posts, "nil")
}

// 获取所有文章
func AllPosts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("per_page"))
	posts, err := models.AllPosts(page, pageSize)
	Response(ctx, 200, &posts, err)
}

// 获取喜欢帖子的用户
func GetFansByPID(ctx *gin.Context) {
	pid := utils.S2uint(ctx.Query("pid"))
	users, err := models.GetUsersByPostID(pid)
	if err != nil {
		Response(ctx, 404, nil, err)
		return
	}
	Response(ctx, 200, users, "Get sucess")
}

// 点赞或取消点赞一个帖子redis
func LikePostByRedis(ctx *gin.Context) {
	uid, pid := ctx.Query("uid"), ctx.Query("pid")
	field := pid + "<-" + uid
	status, _ := cache.IsLike(field)
	cache.LikeOrNot(pid, status, field)
	Response(ctx, 200, nil, "Like sucess")
}
