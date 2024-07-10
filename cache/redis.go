package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx context.Context = context.Background()

const (
	Status     string = "status"
	PLikeCount string = "postLikeCount"
	PLookCount string = "postLookCount"
)

func Init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
	})
	// RDB = redis.NewClient(&redis.Options{
	// 	Addr:     config.Redis["addr"],
	// 	Password: config.Redis["password"],
	// })
}

// 获取浏览量
func LookCount(postID string) (float64, error) {
	lookcount, err := RDB.ZScore(Ctx, PLookCount, postID).Result()
	return lookcount, err
}

// 获取点赞数
func LikeCount(postID string) (uint64, error) {
	likecount, err := RDB.HGet(Ctx, PLikeCount, postID).Uint64()
	return likecount, err
}

// 当前用户是否点赞帖子
func IsLike(field string) (string, error) {
	status, err := RDB.HGet(Ctx, Status, field).Result()
	return status, err
}

// 点赞or取消点赞
func LikeOrNot(postID, status, str string) {
	if status == "1" {
		RDB.HSet(Ctx, Status, str, "0")
		RDB.HIncrBy(Ctx, PLikeCount, postID, -1)
	} else {
		RDB.HSet(Ctx, Status, str, "1")
		RDB.HIncrBy(Ctx, PLikeCount, postID, 1)
	}

}

func SetNX(key string, value interface{}, time time.Duration) (bool, error) {
	exist, err := RDB.SetNX(Ctx, key, value, time).Result()
	return exist, err
}

func Get(key string) (int, error) {
	val, err := RDB.Get(Ctx, key).Int()
	return val, err
}

func Incr(key string) {
	RDB.Incr(Ctx, key)
}

func ZAdd(key string, member string) {
	RDB.ZAdd(Ctx, key, redis.Z{Score: 0, Member: member})
}

func ZRem(key string, member string) {
	RDB.ZRem(Ctx, key, member)
}

func HSet(key string, field string, value interface{}) {
	RDB.HSet(Ctx, key, field, value)
}

func HDel(key string, field string) {
	RDB.HDel(Ctx, key, field)
}

func ZIncrBy(key string, member string) {
	RDB.ZIncrBy(Ctx, key, 1, member)
}

// RDB.ZIncrBy(cache.Ctx, cache.PLookCount, 1, utils.Stringconv(pid))

// // 浏览量， 根评论点赞量， 根评论数量， 根评论内容， 根评论是否点赞

// cache.RDB.RPush(cache.Ctx, "test", 1,2,3,4,5)
// v, _ := cache.RDB.LRange(cache.Ctx, "test", 2, -1).Result()
// cache.RDB.Expire(cache.Ctx, "test", time.Second * 20)
