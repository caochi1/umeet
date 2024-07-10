package models

import (
	"Umeet/cache"
	"Umeet/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Comment struct {
	Model
	CommentID  *uint32    `json:"parentCommentId"`                      //父评论ID *uint表示可以为null
	PostID     *uint32    `json:"pid"`                                  //一对多，帖子对评论
	FavorCount uint32     `json:"favorcount"`                           //点赞数
	UserinfoID uint32     `json:"uid"`                               //用于一对多关系的id
	UserName   string     `json:"authorName" gorm:"size:20"`            //用户昵称
	Content    string     `json:"content"`                              //内容
	Comments   []Comment  `json:"-"`                                    //自引用一对多
	Users      []Userinfo `json:"-" gorm:"many2many:userinfo_comment;"` //与用户的多对多
}

type replyComment struct {
	Comment
	Count  uint32 `json:"commentCount"`
	IsLike bool   `json:"isLike"`
	Avatar string `json:"avatarPath" gorm:"size:255"`
}

// 添加评论
func AddComment(com *Comment, pid uint32) error {
	defer UpdatePost(pid, "comment_count", gorm.Expr("comment_count + ?", 1))
	return Db.Create(com).Error
}

// 删除评论
func DeleteComment(cid, pid uint32) error {
	defer UpdatePost(pid, "comment_count", gorm.Expr("comment_count - ?", 1))
	return Db.Select(clause.Associations).Delete(&Comment{Model: Model{ID: cid}}).Error
}

// 获取根评论
func GetRootComments(pid, uid uint32) []replyComment {
	var comments []Comment
	var replycomments []replyComment
	if err := Db.Model(&Post{Model: Model{ID: pid}}).Association("Comments").Find(&comments); err != nil {
		return nil
	}
	for _, comment := range comments {
		user, _ := GetUserInfo(comment.UserinfoID, "avatar")
		count, status := uint32(GetChildCount(comment.ID)), CommentLikeQuery(comment.ID, uid)
		replycomments = append(replycomments, replyComment{comment, count, status == 1, user.Avatar})
	}
	cache.ZIncrBy(cache.PLookCount, utils.Stringconv(pid))
	return replycomments
}

// 获取子评论
func GetChildComment(cid, uid uint32) []replyComment {
	var comments []Comment
	var replycomments []replyComment
	if err := Db.Model(&Comment{Model: Model{ID: cid}}).Association("Comments").Find(&comments); err != nil {
		return nil
	}
	for _, comment := range comments {
		user, _ := GetUserInfo(comment.UserinfoID, "avatar")
		count, status := uint32(GetChildCount(comment.ID)), CommentLikeQuery(comment.ID, uid)
		replycomments = append(replycomments, replyComment{comment, count, status == 1, user.Avatar})
	}
	return replycomments
}

// 子评论数量
func GetChildCount(cid uint32) int64 {
	return Db.Model(&Comment{Model: Model{ID: cid}}).Association("Comments").Count()
}

// 查询用户是否点赞评论
func CommentLikeQuery(cid, uid uint32) int64 {
	user := Userinfo{Model: Model{ID: uid}}
	count := Db.Model(&user).Where("id = ?", cid).Limit(1).Association("FavorComments").Count()
	return count
}

// 点赞评论
func LikeComment(cid, uid uint32) error {
	user, comment := Userinfo{Model: Model{ID: uid}}, Comment{Model: Model{ID: cid}}
	if err = Db.Model(&user).Association("FavorComments").Append(&comment); err != nil {
		return err
	}
	Db.Model(comment).Update("favor_count", gorm.Expr("favor_count + ?", 1))
	return nil
}

// 取消点赞评论
func DisLikeComment(cid, uid uint32) error {
	user, comment := Userinfo{Model: Model{ID: uid}}, Comment{Model: Model{ID: cid}}
	if err = Db.Model(&user).Association("FavorComments").Delete(&comment); err != nil {
		return err
	}
	Db.Model(comment).Update("favor_count", gorm.Expr("favor_count - ?", 1))
	return nil
}

func AllComment(page, pageSize int) ([]Comment, error) {
	var comments []Comment
	err := Db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&comments).Error
	return comments, err
}
