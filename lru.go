package lru

import (
    "container/list"
    "sync"
    "time"
    "hash/fnv"
)

const NoExpiration int64 = 0

// cache 是 LRU 缓存的核心数据结构
type cache struct {
    mu       sync.Mutex               // 保护并发访问
    list     *list.List               // 双向链表，用于维护 LRU 顺序
    indexMap map[string]*list.Element // 哈希表，用于快速查找节点
    capacity int                      // 缓存容量
}

// 封装的私有节点结构体，仅供 cache 内部使用
type cacheEntry struct {
    key      string
    value    interface{}
    expireAt int64
}

// newCache 创建新的 LRU 缓存实例
func newCache(cap int) *cache {
    return &cache{
        list:     list.New(),
        indexMap: make(map[string]*list.Element),
        capacity: cap,
    }
}

// put 向缓存中插入键值对，支持过期时间
func (c *cache) put(k string, v interface{}, expireAt int64) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // 如果键已经存在，更新值并将该节点移到链表头部
    if elem, exists := c.indexMap[k]; exists {
        entry := elem.Value.(*cacheEntry)
        entry.value = v
        entry.expireAt = expireAt
        c.list.MoveToFront(elem) // 最近访问的移到头部
        return
    }

    // 如果缓存已满，移除最久未使用的节点（链表末尾）
    if c.list.Len() >= c.capacity {
        oldest := c.list.Back()
        if oldest != nil {
            oldestEntry := oldest.Value.(*cacheEntry)
            delete(c.indexMap, oldestEntry.key) // 删除索引
            c.list.Remove(oldest)               // 移除链表中的节点
        }
    }

    // 插入新节点到链表头部
    entry := &cacheEntry{key: k, value: v, expireAt: expireAt}
    elem := c.list.PushFront(entry)
    c.indexMap[k] = elem
}

// get 从缓存中获取键对应的值，如果存在且未过期则返回值，并将该节点移到链表头部
func (c *cache) get(k string) (interface{}, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if elem, exists := c.indexMap[k]; exists {
        entry := elem.Value.(*cacheEntry)
        // 如果不过期或者当前时间还未超过过期时间，则返回值
        if entry.expireAt == NoExpiration || time.Now().UnixNano() < entry.expireAt {
            c.list.MoveToFront(elem) // 访问后移到链表头部
            return entry.value, true
        }

        // 如果已经过期，删除该节点
        delete(c.indexMap, k)
        c.list.Remove(elem)
    }
    return nil, false
}

// del 从缓存中删除指定键的条目
func (c *cache) del(k string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if elem, exists := c.indexMap[k]; exists {
        delete(c.indexMap, k)
        c.list.Remove(elem)
    }
}

// Cache 结构体表示整个 LRU/LRU-2 缓存系统，支持分片
type Cache struct {
    shards     []*cache      // 分片缓存，用于减少锁争用
    shardCount int           // 分片数量
    expiration time.Duration // 默认过期时间
}

// NewLRUCache 创建一个支持分片和过期时间的 LRU 缓存系统
func NewLRUCache(bucketCnt, capPerBkt int, expiration ...time.Duration) *Cache {
    c := &Cache{
        shards:     make([]*cache, bucketCnt),
        shardCount: bucketCnt,
        expiration: 0,
    }

    // 如果指定了过期时间，则设置
    if len(expiration) > 0 {
        c.expiration = expiration[0]
    }

    // 初始化每个分片
    for i := 0; i < bucketCnt; i++ {
        c.shards[i] = newCache(capPerBkt)
    }

    return c
}

// hash 使用 FNV 哈希算法计算键的哈希值
func hash(s string) uint32 {
    h := fnv.New32a()
    _, _ = h.Write([]byte(s))
    return h.Sum32()
}

// getShard 根据键获取对应的缓存分片
func (c *Cache) getShard(key string) *cache {
    h := hash(key)
    return c.shards[h%uint32(c.shardCount)]
}

// Put 向缓存中插入一个条目，支持过期时间
func (c *Cache) Put(key string, val interface{}) {
    shard := c.getShard(key)
    expireAt := NoExpiration
    if c.expiration > 0 {
        expireAt = time.Now().Add(c.expiration).UnixNano() // 计算过期时间
    }
    shard.put(key, val, expireAt)
}

// Get 从缓存中获取一个条目
func (c *Cache) Get(key string) (interface{}, bool) {
    shard := c.getShard(key)
    return shard.get(key)
}

// Del 删除缓存中的一个条目
func (c *Cache) Del(key string) {
    shard := c.getShard(key)
    shard.del(key)
}

// LRU2 启用 LRU-2 缓存，创建第二级缓存
func (c *Cache) LRU2(capPerBkt int) *Cache {
    for i := range c.shards {
        c.shards[i] = newCache(capPerBkt) // 创建第二级缓存
    }
    return c
}
