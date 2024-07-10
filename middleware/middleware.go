package middleware

import (
	"Umeet/cache"
	"Umeet/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// JWT鉴权
func Authorization(ctx *gin.Context) {
	auth := ctx.GetHeader("authorization")
	token := strings.Split(auth, " ")
	if auth == "" || len(token) != 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	_, err := utils.ParseToken(token[1])
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Next()
}

// func checkToken(claims utils.UserClaims) bool {
// 	time := time.Now().Unix() - claims.ExpiresAt.Unix()
// 	var week int64 = 60 * 60 * 24 * 7
// 	return time < week
// }

// ip限流
func IPLimiter(ctx *gin.Context) {

	ip := ctx.ClientIP()
	if exist, _ := cache.SetNX(ip, 1, time.Second*10); !exist {
		cache.Incr(ip)
	}
	if frequency, _ := cache.Get(ip); frequency >= 50 {
		ctx.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	// limiter := rate.NewLimiter(10, 1)
	// list.List
	ctx.Next()
}
