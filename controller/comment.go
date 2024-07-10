package controller

import (
	"Umeet/models"
	"Umeet/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 添加根评论
func AddRootComment(ctx *gin.Context) {
	var comment models.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		Response(ctx, 404, nil, err)
		return
	}
	if err := models.AddComment(&comment, *comment.PostID); err != nil {
		Response(ctx, 404, nil, err)
		return
	}
	Response(ctx, 200, comment, "AddNewComment sucess")
}

// 添加子评论
func AddChildComment(ctx *gin.Context) {
	var comment models.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		Response(ctx, 404, nil, err)
		return
	}
	pid := comment.PostID
	comment.PostID = nil
	if err := models.AddComment(&comment, *pid); err != nil {
		Response(ctx, 404, nil, err)
		return
	}
	Response(ctx, 200, comment, "AddNewComment sucess")
}

// 删除评论
func DeleteComment(ctx *gin.Context) {
	cid, pid := utils.S2uint(ctx.Query("cid")), utils.S2uint(ctx.Query("pid"))
	if err := models.DeleteComment(cid, pid); err != nil {
		Response(ctx, 404, nil, err)
		return
	}
	Response(ctx, 200, cid, "DeleteComment sucess")
}

// 获取帖子的根评论
func GetRootComments(ctx *gin.Context) {
	uid, pid := utils.S2uint(ctx.Query("uid")), utils.S2uint(ctx.Query("pid"))
	comments := models.GetRootComments(pid, uid)
	Response(ctx, 200, comments, "GetRootComments sucess")
}

// 获取子评论
func GetChildComment(ctx *gin.Context) {
	uid, cid := utils.S2uint(ctx.Query("uid")), utils.S2uint(ctx.Query("cid"))
	comments := models.GetChildComment(cid, uid)
	Response(ctx, 200, comments, "GetChildComments sucess")
}

// 点赞评论
func LikeComment(ctx *gin.Context) {
	uid, cid := utils.S2uint(ctx.Query("uid")), utils.S2uint(ctx.Query("cid"))
	status := ctx.Query("islike")
	if status == "true" {
		models.DisLikeComment(cid, uid)
	} else {
		models.LikeComment(cid, uid)
	}
	Response(ctx, 200, nil, nil)
}

// 所有评论
func AllComment(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("per_page"))
	comments, err := models.AllComment(page, pageSize)
	Response(ctx, 200, &comments, err)
}
