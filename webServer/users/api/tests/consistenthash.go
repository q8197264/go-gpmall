package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

/**
	一致性哈希实现 demo
**/
type Hash func([]byte) uint32

type Map struct {
	replies int   //虚拟节点倍数
	keys    []int //哈希环，已排序
	hash    Hash  //哈希函数，默认crc32

	hashMap map[int]string //虚拟节点与真实节点映射关系
}

func main1() {
	m := NewConsistenthash(3, nil)
	m.Add("192.168.1.10:800", "192.168.1.11:801", "192.168.1.112:802")
	m.Get("192.168.1.110:800")
}

func NewConsistenthash(replies int, hash Hash) *Map {
	m := &Map{
		replies: replies,
		hash:    hash,
		hashMap: make(map[int]string),
	}
	if hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// 添加节点: 可一次性添加多个服务器
func (m *Map) Add(keys ...string) {
	for _, key := range keys {

		// 虚拟节点的名称为序号+key
		for i := 0; i < m.replies; i++ {

			// 计算虚拟节点的哈希值，并进行类转换
			virNodeHash := int(m.hash([]byte(key + strconv.Itoa(i))))

			// 把虚拟节点加入hash环
			m.keys = append(m.keys, virNodeHash)

			// 维护虚拟节点到真实节点之间的映射关系
			m.hashMap[virNodeHash] = key
		}
	}

	// 对哈希环进行排序
	sort.Ints(m.keys)
}

// Get 选择存储节点的函数
func (m *Map) Get(key string) string {

	// ------ debug begin -----
	fmt.Printf("%v\n", m.keys)
	for k, v := range m.hashMap {
		println(k, v)
	}
	// ------- debug end --------

	if key == "" {
		return ""
	}

	// 1. 计算key的哈希值
	keyhash := int(m.hash([]byte(key)))

	// 2. 顺时针找到哈希环上 第一个 匹配的 相邻虚拟节点 的 下标
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= keyhash
	})

	// 3. 通过 下标 定位到虚拟节点
	virNodeHash := m.keys[index%len(m.keys)]

	// 4. 从虚拟节点打到真实节点
	node := m.hashMap[virNodeHash]

	println("=", virNodeHash, node)

	return node
}
