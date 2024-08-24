package cache

import (
	"errors"
	"fmt"
	"github.com/bluele/gcache"
	"sync"
	"time"
)

var once sync.Once
var Cache *cache

type cache struct {
	cache gcache.Cache //该对象可以缓存任何类型数据
}

const (
	cacheSize = 1_000_000
	cacheTTL  = 2 * time.Hour // default expiration
)

var (
	errUserNotInCache = errors.New("the user isn't in cache")
)

func Init() *cache {
	once.Do(func() {
		Cache = newGCache()
	})
	return Cache
}

// NewGCache 创建缓存对象，使用ARC算法淘汰缓存元素
func newGCache() *cache {
	return &cache{
		cache: gcache.New(cacheSize).LRU().Build(),
	}
}

// 更新缓存元素过期时间
func (gc *cache) Update(u string, expireIn time.Duration) error {
	fmt.Println("更新缓存", u, expireIn)
	return gc.cache.SetWithExpire(u, u, expireIn)
}

// 读取缓存
func (gc *cache) Read(key string) (string, error) {
	val, err := gc.cache.Get(key)
	if err != nil {
		if errors.Is(err, gcache.KeyNotFoundError) {
			return "", errUserNotInCache
		}
		return "", fmt.Errorf("get: %w", err)
	}
	return val.(string), nil
}

// 更新缓存元素过期时间
func (gc *cache) SetCache(key string, value any, expireIn time.Duration) error {
	return gc.cache.SetWithExpire(key, value, expireIn)
}

// 读取缓存
func (gc *cache) GetCache(key string) (any, error) {
	val, err := gc.cache.Get(key)
	if err != nil {
		if errors.Is(err, gcache.KeyNotFoundError) {
			return nil, errUserNotInCache
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	gc.cache.Remove(key)
	return val, nil
}

// 弹出一个值
func (gc *cache) PullCache(key string) (any, error) {
	val, err := gc.cache.Get(key)
	if err != nil {
		if errors.Is(err, gcache.KeyNotFoundError) {
			return nil, errUserNotInCache
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	gc.cache.Remove(key)
	return val, nil
}

// Delete 删除缓存元素
func (gc *cache) Delete(key string) {
	fmt.Println("删除缓存", key)
	gc.cache.Remove(key)
}
