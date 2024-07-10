package models

type Image struct {
	ID     uint32 `gorm:"primarykey"`
	Dst    string `gorm:"size:255"`
	PostID uint32 
}

// 保存图片路径
func SavePostPath(dst string, pid uint32) error {
	var image Image
	image.Dst, image.PostID = dst, pid
	err := Db.Create(&image).Error
	return err
}

// 获取图片路径
func GetPostImages(pid uint32) ([]Image, error) {
	var images []Image
	err := Db.Select("dst").Model(&Post{Model: Model{ID: pid}}).Association("Images").Find(&images)
	return images, err
}
