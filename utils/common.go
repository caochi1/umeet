package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"sync"
)

func Stringconv(id uint32) string { return strconv.Itoa(int(id)) }

func EncryMd5(s string) string {
	ctx := md5.New()
	ctx.Write([]byte(s))
	return hex.EncodeToString(ctx.Sum(nil))
}
func S2uint(id string) uint32 {
	uintid, _ := strconv.Atoi(id)
	return uint32(uintid)
}

// 高性能字符串拼接
func StringsBuilder(str ...string) string {
	var builder strings.Builder
	for i := 0; i < len(str); i++ {
		builder.WriteString(str[i])
	}
	return builder.String()
}

// ########################################################
type SafeMap struct {
	m    map[any]any
	lock sync.RWMutex
}

func NewSafeMap(len uint8) *SafeMap {
	return &SafeMap{
		make(map[any]any, len),
		sync.RWMutex{},
	}
}

func (sm *SafeMap) Set(key, value any) {
	sm.lock.Lock()
	sm.m[key] = value
	sm.lock.Unlock()
}

func (sm *SafeMap) Get(key any) (any, bool) {
	sm.lock.RLock()
	value, ok := sm.m[key]
	sm.lock.RUnlock()
	return value, ok
}

func (sm *SafeMap) Delete(key any) {
	sm.lock.Lock()
	delete(sm.m, key)
	sm.lock.Unlock()
}

func (sm *SafeMap) Len() int {
	return len(sm.m)
}

func (sm *SafeMap) ForEach(f func(k, v any)) {
	sm.lock.RLock()
	for k, v := range sm.m {
		f(k, v)
	}
	sm.lock.RUnlock()
}

// func (sm *SafeMap) locker(f func(...any)) func(...any) {
// 	return func(args ...any) {
// 		if len(args) == 2 {

// 		}
// 		sm.lock.Lock()
// 		f(client)
// 		sm.lock.Unlock()
// 	}
// }

// func insert(arr []int, target, idx int) []int {
// 	b := make([]int, len(arr))
// 	copy(b, arr)
// 	arr = append(append(b[:idx], target), arr[idx:]...)
// 	return arr
// }

// func binary_search(arr []int, target, l, r int) int {
// 	for l <= r {
// 		middle := (l + r) / 2
// 		if arr[middle] == target {
// 			return middle
// 		} else if arr[middle] < target {
// 			l = middle + 1
// 		} else {
// 			r = middle - 1
// 		}
// 	}
// 	return l
// }
