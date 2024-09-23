# Go 高性能 LRU/LRU-2 缓存库

<p align="center">
  <a href="/go.mod#L3" alt="go version">
    <img src="https://img.shields.io/badge/go%20version-%3E=1.18-brightgreen?style=flat"/>
  </a>
  <a href="https://goreportcard.com/badge/github.com/orca-zhang/ecache" alt="goreport">
    <img src="https://goreportcard.com/badge/github.com/orca-zhang/ecache">
  </a>
  <a href="https://orca-zhang.semaphoreci.com/projects/ecache" alt="buiding status">
    <img src="https://orca-zhang.semaphoreci.com/badges/ecache.svg?style=shields">
  </a>
  <a href="https://codecov.io/gh/orca-zhang/ecache" alt="codecov">
    <img src="https://codecov.io/gh/orca-zhang/ecache/branch/master/graph/badge.svg?token=F6LQbADKkq"/>
  </a>
  <a href="https://github.com/orca-zhang/ecache/blob/master/LICENSE" alt="license MIT">
    <img src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat">
  </a>
</p>

## 项目简介

**Go LRU/LRU-2 缓存库** 是一个高效的、并发安全的缓存解决方案，适用于处理高频访问的热点数据。该库基于 [Least Recently Used (LRU)](https://en.wikipedia.org/wiki/Cache_replacement_policies#LRU) 和 LRU-2 算法实现，提供了先进的数据存储和管理策略，能在性能、内存占用和数据准确性之间取得良好平衡。

通过该缓存库，你可以：
- **快速存取** 高频访问的数据。
- **自动移除** 长时间未使用的数据，释放内存。
- **过期管理**：支持数据条目的过期时间控制。
- **支持并发**：使用分片技术，避免并发冲突，提升性能。
- **LRU-2 算法**：提高热点数据的访问命中率，适合需要高效处理访问频次的场景。

## 特性

- **O(1) 复杂度**：基于双向链表和哈希表的设计，保证插入、查找、删除操作的时间复杂度为 O(1)。
- **分片缓存**：通过分片设计减少锁竞争，提升多线程并发场景下的性能。
- **数据过期**：支持设置条目过期时间，自动清除过期数据，节省内存。
- **LRU-2 支持**：访问两次后将数据条目提升至高级缓存，提高热点数据命中率。

## 安装

使用 Go 模块进行安装：

```bash
go get github.com/7836246/lru
```

确保你的 Go 版本在 `1.18` 或以上。

## 使用方法

下面是详细的使用示例，展示了如何创建缓存、插入和读取数据、处理过期数据，以及使用 LRU-2 缓存：

### 基础使用

```go
package main

import (
    "fmt"
    "time"
    "github.com/7836246/lru"
)

func main() {
    // 创建一个带有 16 个分片、每个分片 128 个条目容量，过期时间为 1 秒的 LRU 缓存
    cache := lru.NewLRUCache(16, 128, 1*time.Second)

    // 插入缓存条目
    cache.Put("key1", "value1")

    // 获取缓存条目
    if val, found := cache.Get("key1"); found {
        fmt.Println("获取到:", val) // 输出: 获取到: value1
    } else {
        fmt.Println("未找到 key1")
    }

    // 测试过期条目
    time.Sleep(2 * time.Second)
    if _, found := cache.Get("key1"); !found {
        fmt.Println("key1 已过期")
    }

    // 测试 LRU-2
    cache.LRU2(128)
    cache.Put("key2", "value2")
    if val, found := cache.Get("key2"); found {
        fmt.Println("LRU-2 缓存获取到:", val) // 输出: LRU-2 缓存获取到: value2
    }
}
```

## 使用场景

1. **高并发访问**：在 Web 服务、缓存代理等需要处理大量并发请求的场景下，可以通过该缓存库来加速数据访问，减少数据库的压力。
2. **热点数据缓存**：适用于频繁访问的热点数据缓存，减少频繁的数据加载和反复计算。
3. **临时数据存储**：可以存储会话数据、短期内频繁访问的数据，并设置自动过期时间。

## 运行测试

运行单元测试来验证缓存库的功能：

```bash
go test -v ./...
```

## 贡献

欢迎提交 issue 和 pull request 来改进这个项目。
