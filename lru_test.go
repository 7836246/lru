package lru_test

import (
    "testing"
    "time"
    "github.com/7836246/lru"
)

func TestLRUCache(t *testing.T) {
    cache := lru.NewLRUCache(16, 128, 1*time.Second)

    cache.Put("key1", "value1")
    if val, found := cache.Get("key1"); !found || val != "value1" {
        t.Errorf("预期得到 value1，实际得到 %v", val)
    }

    // 测试过期
    time.Sleep(2 * time.Second)
    if _, found := cache.Get("key1"); found {
        t.Errorf("预期 key1 已过期")
    }

    // 测试删除
    cache.Put("key2", "value2")
    cache.Del("key2")
    if _, found := cache.Get("key2"); found {
        t.Errorf("预期 key2 已被删除")
    }
}

func TestLRU2Cache(t *testing.T) {
    cache := lru.NewLRUCache(16, 128).LRU2(64)

    cache.Put("key1", "value1")
    cache.Put("key2", "value2")

    // 第一次访问 key1
    if val, found := cache.Get("key1"); !found || val != "value1" {
        t.Errorf("第一次预期得到 value1，实际得到 %v", val)
    }

    // 第二次访问 key1
    if val, found := cache.Get("key1"); !found || val != "value1" {
        t.Errorf("第二次预期得到 value1，实际得到 %v", val)
    }
}
