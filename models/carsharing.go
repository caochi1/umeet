package models

import (
	"Umeet/utils"
	"Umeet/websocket"
	"strconv"
	"time"
)

type CarSharing struct {
	ID          uint32 `gorm:"primarykey"`
	UID         uint32 `json:"id"`
	Content     string `json:"content" gorm:"size:100"`
	StartPoint  string `json:"departure" gorm:"size:50"`
	EndPoint    string `json:"destination" gorm:"size:50"`
	Population  uint8  `json:"population"`
	ExpiresAt   int64  `json:"time"`
	CreatedAt   time.Time
	// ChatRecords []ChatRecords `json:"-"`
}

// type ChatRecords struct {
// 	ID           uint32 `gorm:"primarykey"`
// 	UID          uint32
// 	Content      string `json:"content" gorm:"size:200"`
// 	CreatedAt    time.Time
// 	CarSharingID uint32
// }

type carReply struct {
	Members map[string]string `json:"members"`
	*CarSharing
}

// 添加帖子
func AddCarPost(cs *CarSharing) error {
	if err := Db.Create(cs).Error; err == nil {
		clientManager := websocket.NewManager()
		go clientManager.Run()
		websocket.Rooms.Set(utils.Stringconv(cs.ID), clientManager)
	}
	return err

}

// 所有拼车帖子
func AllCarPosts() ([]carReply, error) {
	var carPosts []carReply
	if websocket.Rooms.Len() == 0 {
		return carPosts, nil
	}
	rows, err := Db.Model(&CarSharing{}).Rows()
	if err != nil {
		return carPosts, err
	}
	for rows.Next() {
		var cs CarSharing
		Db.ScanRows(rows, &cs)
		cm := websocket.Rooms.Get(utils.Stringconv(cs.ID))

		if time.Now().Unix() < cs.ExpiresAt {
			carPosts = append(carPosts, carReply{cm.Copy(), &cs})
		} else {
			Db.Delete(&CarSharing{}, cs.ID)
			cm.Close()
			websocket.Rooms.Delete(strconv.Itoa(int(cs.ID)))
		}
	}
	rows.Close()
	return carPosts, err
}
