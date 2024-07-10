package cache

import (
	"Umeet/config"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestSomething(t *testing.T) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Redis["addr"],
		Password: "",
		DB:       0,
	})
	// var cursor uint64
	// keys, _, _ := RDB.Scan(Ctx, cursor, "*", 100).Result()
	// for _, key := range keys {
	// 	fmt.Println("key:", key)
	// }
}
