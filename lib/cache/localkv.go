package cache

import (
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

type localCacheValue struct {
	expiredAt int64
	content   any
}

type localCache struct {
	mux  *sync.Mutex
	maps map[string]*localCacheValue
}

var (
	lc   *localCache
	once sync.Once
)

func init() {
	once.Do(func() {
		lc = &localCache{
			mux:  &sync.Mutex{},
			maps: make(map[string]*localCacheValue),
		}
	})
}

// Put stores a key-value pair
//
// The default expiration is one hour if no expiration is specified
func Put(key string, val any, expiredTimeDuration ...time.Duration) bool {
	var expiredAt int64
	if len(expiredTimeDuration) > 0 {
		expiredAt = utils.NewTimeWithTimeZone().Now().Add(expiredTimeDuration[0]).Unix()
	} else {
		expiredAt = utils.NewTimeWithTimeZone().Now().Add(time.Hour).Unix()
	}
	lc.mux.Lock()
	lc.maps[key] = &localCacheValue{expiredAt: expiredAt, content: val}
	lc.mux.Unlock()

	return true
}

// Get retrieves the value for a key
// ret indicates (0:key exists and has not expired,-1:key has expired,-2:key does not exist)
func Get(key string) (value any, ret int) {
	lc.mux.Lock()
	defer lc.mux.Unlock()
	if val, ok := lc.maps[key]; ok {
		if val.expiredAt < utils.NewTimeWithTimeZone().Now().Unix() {
			delete(lc.maps, key) // <- Delete expired keys
			return nil, -1
		}
		value = val.content
		ret = 0
		return
	} else {
		return nil, -2
	}
}

// Forget deletes the value for a key
func Forget(key string) bool {
	lc.mux.Lock()
	defer lc.mux.Unlock()
	if _, ok := lc.maps[key]; ok {
		delete(lc.maps, key)
		return true
	}

	return false
}

// Expire sets the expiration for a key
func Expire(key string, expiredTimeDuration time.Duration) bool {
	var optOk bool

	expiredAt := utils.NewTimeWithTimeZone().Now().Add(expiredTimeDuration).Unix()

	lc.mux.Lock()
	if _, ok := lc.maps[key]; ok {
		lc.maps[key].expiredAt = expiredAt
		optOk = true
	}
	lc.mux.Unlock()

	return optOk
}
