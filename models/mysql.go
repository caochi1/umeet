package models

import (
	"Umeet/cache"
	"Umeet/config"
	"Umeet/utils"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	err error
	Db  *gorm.DB
)

type Model struct {
	ID        uint32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Init(migrate bool) {
	// DataSourceName := "root:20011025Hp-@tcp(127.0.0.1:3306)/umeet?charset=utf8mb4&parseTime=True&loc=Local"
	DataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.MySQL["username"],
		config.MySQL["password"],
		config.MySQL["host"],
		config.MySQL["port"],
		config.MySQL["dbname"],
		config.MySQL["charset"])
	Db, err = gorm.Open(mysql.Open(DataSourceName))
	// Db, err = gorm.Open(mysql.Open(DataSourceName), &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Info),
	// })
	if err != nil || Db.Error != nil {
		log.Println("MySQL open failed")
		return
	}
	if migrate {
		// err = Db.AutoMigrate(&CarSharing{})
		err = Db.AutoMigrate(&Userinfo{}, &Post{}, &Comment{}, &Image{}, &CarSharing{})
		if err != nil {
			log.Println("AutoMigrate failed")
		}
	}
	// cronJob()
}

// 将redis数据迁移至mysql
func Rtom() {
	//点赞量
	go func() {
		for _, pid := range cache.RDB.HKeys(cache.Ctx, cache.PLikeCount).Val() {
			count, _ := cache.RDB.HGet(cache.Ctx, cache.PLikeCount, pid).Uint64()
			pid := utils.S2uint(pid)
			if CheckPost(pid, &Post{}) == 0 {
				continue
			}
			UpdatePost(pid, "favor_count", uint32(count))
		}
	}()
	//浏览量
	go func() {
		for _, v := range cache.RDB.ZRangeWithScores(cache.Ctx, cache.PLookCount, 0, -1).Val() {
			if pid, ok := (v.Member).(string); ok {
				count := v.Score
				UpdatePost(utils.S2uint(pid), "look_count", uint32(count))
			}
		}
	}()
	//点赞状态
	go func() {
		for relation, status := range cache.RDB.HGetAll(cache.Ctx, cache.Status).Val() {
			id := strings.Split(relation, "<-")
			uid, pid := utils.S2uint(id[1]), utils.S2uint(id[0])
			switch {
			case CheckPost(pid, &Post{}) == 0:
			case status == "1":
				LikePost(pid, uid)
			default:
				DisLikePost(pid, uid)
			}
		}
	}()

	cache.RDB.FlushAll(cache.Ctx)
}

// 定时任务
func cronJob() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("@every 24h", Rtom)
	c.Start()
}
