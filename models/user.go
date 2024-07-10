package models

type Userinfo struct {
	Model
	Description   string    `json:"description" gorm:"size:50"`              //简介
	Gender        string    `json:"gender" gorm:"size:20"`                   //性别
	NickName      string    `json:"nickname" gorm:"size:40" `                //昵称
	UserName      string    `json:"username" gorm:"index:idx_name size:40" ` //用户账号
	Passwd        string    `json:"password" gorm:"size:40" `                //密码
	Avatar        string    `json:"avatarPath" gorm:"size:255"`              //头像
	Faculty       string    `json:"faculty" gorm:"size:80"`                  //学院
	Program       string    `json:"program" gorm:"size:80"`                  //学位
	PostList      []Post    `json:"-"`                                       //发布的帖子
	FavorPosts    []Post    `json:"-" gorm:"many2many:userinfo_post;"`       //喜欢的帖子
	Comments      []Comment `json:"-"`                                       //发布的评论
	FavorComments []Comment `json:"-" gorm:"many2many:userinfo_comment;"`    //喜欢的评论
}

// *string 可以为空值 NULL

// 查询用户是否存在+取值
func CheckUser(username string) (Userinfo, error) {
	var user Userinfo
	err := Db.
		Select("id", "gender", "nick_name", "faculty", "program", "passwd", "avatar").
		Where("user_name = ?", username).Take(&user).Error
	return user, err
}

// 查询用户是否存在
func VerifyUser(username string) int64 {
	return Db.Select("id").Model(Userinfo{}).Where("user_name = ?", username).Limit(1).RowsAffected
}

// 创建用户
func CreateUser(user *Userinfo) error {
	err := Db.Create(user).Error
	return err
}

// 删除用户
func RmUser(uid uint32) {
	Db.Delete(&Userinfo{}, uid)
}
func AllUsers(page, pageSize int) ([]Userinfo, error) {
	var users []Userinfo
	err := Db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&users).Error
	return users, err
}

// 获取用户信息
func GetUserInfo(uid uint32, column interface{}) (Userinfo, error) {
	var user Userinfo
	err := Db.Select(column).Where("id = ?", uid).Take(&user).Error
	return user, err
}

// 更新用户信息many
func UpdatesUser(uid uint32, data Userinfo) error {
	return Db.Model(&Userinfo{}).Omit("id").Where("id = ?", uid).Updates(data).Error
}

// 更新用户信息single
func UpdateUser(uid uint32, column string, data interface{}) error {
	return Db.Model(&Userinfo{}).Where("id = ?", uid).Update(column, data).Error
}

// 获取用户发布的帖子
func GetPostByUserID(uid uint32) ([]replyPost, error) {
	var posts []Post
	var replyposts []replyPost
	var user Userinfo
	err = Db.Model(&Userinfo{Model: Model{ID: uid}}).Association("PostList").Find(&posts)
	Db.Model(&Userinfo{}).Select("nick_name", "avatar").Take(&user, uid)
	for _, post := range posts {
		post.FavorCount, post.LookCount = PCount(post.ID)
		rp := replyPost{post, IsLike(uid, post.ID) == "1", user.NickName, user.Avatar}
		replyposts = append(replyposts, rp)
	}
	return replyposts, err
}

// 获取用户喜欢的帖子
func GetFavorPostsByUserID(uid uint32) ([]Post, error) {
	var user Userinfo
	err := Db.Preload("FavorPosts").Take(&user, uid).Error
	return user.FavorPosts, err
}
