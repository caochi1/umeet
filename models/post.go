package models

import (
	"Umeet/cache"
	"Umeet/utils"
	"os"
	"strings"

	"gorm.io/gorm/clause"
)

type Post struct {
	Model
	Title        string     `json:"title" gorm:"size:100"`             //标题
	Content      string     `json:"content" gorm:"type:TEXT"`          //内容
	LookCount    uint32     `json:"lookcount"`                         //浏览量
	CommentCount uint32     `json:"commentcount"`                      //评论数
	FavorCount   uint32     `json:"favorcount"`                        //点赞量
	UserinfoID   uint32     `json:"authorId"`                          //发布帖子的用户ID
	Users        []Userinfo `json:"-" gorm:"many2many:userinfo_post;"` //喜欢帖子的用户
	Comments     []Comment  `json:"-"`                                 //拥有的评论
	Images       []Image    `json:"-"`                                 //拥有的图片
}

type replyPost struct {
	Post
	IsLike   bool   `json:"isLike"`                     //是否点赞
	UserName string `json:"authorname" gorm:"size:20"`  //发布帖子的用户名
	Avatar   string `json:"avatarPath" gorm:"size:255"` //用户头像
}

// gorm:"default:0"

// 检查帖子是否存在
func CheckPost(id uint32, model interface{}) int64 {
	return Db.Select("id").Take(model, id).RowsAffected
}

// 添加帖子
func AddPost(post *Post) (uint32, error) {
	if err := Db.Create(post).Error; err == nil {
		cache.ZAdd(cache.PLookCount, utils.Stringconv(post.ID))
	}
	return post.ID, err
}

// 删除帖子
func DeletePost(pid uint32) error {
	images, _ := GetPostImages(pid)
	if err := Db.Select(clause.Associations).Delete(&Post{Model: Model{ID: pid}}).Error; err == nil {
		id := utils.Stringconv(pid)
		cache.HDel(cache.PLikeCount, id)
		cache.ZRem(cache.PLookCount, id)
		for _, v := range images {
			dst := strings.Split(v.Dst, "/")
			os.Remove(utils.StringsBuilder("static/post/", dst[len(dst)-1]))
		}
	}
	return err
}

// 点赞帖子
func LikePost(pid, uid uint32) error {
	user, post := Userinfo{Model: Model{ID: uid}}, Post{Model: Model{ID: pid}}
	if err = Db.Model(&user).Association("FavorPosts").Append(&post); err != nil {
		return err
	}
	return nil
}

// 取消点赞帖子
func DisLikePost(pid, uid uint32) error {
	user, post := Userinfo{Model: Model{ID: uid}}, Post{Model: Model{ID: pid}}
	if err = Db.Model(&user).Association("FavorPosts").Delete(&post); err != nil {
		return err
	}
	return nil
}

// 查询用户是否点赞帖子
func PostLikeQuery(pid, uid uint32) int64 {
	user := Userinfo{Model: Model{ID: uid}}
	count := Db.Model(&user).Where("id = ?", pid).Limit(1).Association("FavorPosts").Count()
	return count
}

// 获取喜欢帖子的所有用户
func GetUsersByPostID(pid uint32) ([]Userinfo, error) {
	var users []Userinfo
	err := Db.Model(&Post{Model: Model{ID: pid}}).Association("Users").Find(&users)
	return users, err
}

// 更改帖子many
func UpdatesPost(pid uint32, data map[string]interface{}) error {
	return Db.Model(&Post{}).Where("id = ?", pid).Updates(data).Error
}

// 更改帖子single
func UpdatePost(pid uint32, column string, value interface{}) error {
	return Db.Model(&Post{}).Where("id = ?", pid).Update(column, value).Error
}

// 获取所有帖子
func AllPosts(page, pageSize int) ([]Post, error) {
	var posts []Post
	err := Db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&posts).Error
	return posts, err
}

// 随机获取主页帖子
func GetPostsByRandom(uid uint32) []replyPost {
	var posts []replyPost
	// rows, err := Db.Model(&Post{}).Limit(6).Order("rand()").Rows()
	rows, err := Db.Model(&Post{}).Limit(6).Order("id desc").Rows()
	if err != nil {
		return nil
	}
	for rows.Next() {
		var post Post
		var user Userinfo
		Db.Select("id", "created_at", "title", "content", "comment_count", "userinfo_id").ScanRows(rows, &post)
		Db.Model(&Userinfo{}).Select("nick_name", "avatar").Where("id = ?", post.UserinfoID).Take(&user)
		post.FavorCount, post.LookCount = PCount(post.ID)
		posts = append(posts, replyPost{post, IsLike(uid, post.ID) == "1", user.NickName, user.Avatar})
	}

	rows.Close()
	return posts
}

// 查询点赞状态
func IsLike(uid, pid uint32) string {
	field := utils.StringsBuilder(utils.Stringconv(pid), "<-", utils.Stringconv(uid))
	status, err := cache.IsLike(field)
	if err != nil {
		if status := PostLikeQuery(pid, uid); status == 0 {
			cache.HSet("status", field, "0")
			return "0"
		} else {
			cache.HSet("status", field, "1")
			return "1"
		}
	}
	return status
}

// 查询点赞和浏览量
func PCount(pid uint32) (uint32, uint32) {
	field := utils.Stringconv(pid)
	likecount, err := cache.LikeCount(field)
	lookcount, _ := cache.LookCount(field)
	if err != nil {
		var p Post
		Db.Where("id = ?", pid).Select("favor_count").Take(&p)
		cache.HSet("postLikeCount", field, uint64(p.FavorCount))
		return p.FavorCount, uint32(lookcount)
	}
	return uint32(likecount), uint32(lookcount)
}
