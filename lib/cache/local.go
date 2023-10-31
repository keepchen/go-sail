package cache

import (
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

type value struct {
	expiredAt int64
	content   interface{}
}

type localCache struct {
	mux  *sync.Mutex
	maps map[string]*value
}

var lc *localCache

func init() {
	(&sync.Once{}).Do(func() {
		lc = &localCache{
			mux:  &sync.Mutex{},
			maps: make(map[string]*value),
		}
	})
}

// Put 保存key-value键值对
//
// 若不设置过期时间，则默认过期时间为一个小时
func Put(key string, val interface{}, expiredTimeDuration ...time.Duration) bool {
	var expiredAt int64
	if len(expiredTimeDuration) > 0 {
		expiredAt = utils.NewTimeWithTimeZone().Now().Add(expiredTimeDuration[0]).Unix()
	} else {
		expiredAt = utils.NewTimeWithTimeZone().Now().Add(time.Hour).Unix()
	}
	lc.mux.Lock()
	lc.maps[key] = &value{expiredAt: expiredAt, content: val}
	lc.mux.Unlock()

	return true
}

// Get 根据key获取对应的value值
// ret的表示(0:key存在且未过期,-1:key已过期,-2:key不存在)
func Get(key string) (value interface{}, ret int) {
	lc.mux.Lock()
	defer lc.mux.Unlock()
	if val, ok := lc.maps[key]; ok {
		if val.expiredAt < utils.NewTimeWithTimeZone().Now().Unix() {
			delete(lc.maps, key) // <- 删除过期key
			return nil, -1
		}
		value = val.content
		ret = 0
		return
	} else {
		return nil, -2
	}
}

// Forget 根据key删除对应的value值
func Forget(key string) bool {
	lc.mux.Lock()
	defer lc.mux.Unlock()
	if _, ok := lc.maps[key]; ok {
		delete(lc.maps, key)
		return true
	}

	return false
}

// Expire 为key设置过期时间
func Expire(key string, expiredTimeDuration time.Duration) bool {
	expiredAt := utils.NewTimeWithTimeZone().Now().Add(expiredTimeDuration).Unix()

	lc.mux.Lock()
	lc.maps[key].expiredAt = expiredAt
	lc.mux.Unlock()

	return true
}
