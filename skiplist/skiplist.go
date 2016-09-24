package skiplist

import (
	"errors"
	"fmt"

	"sync/atomic"

	"github.com/supermeng/godatastruct/common"
)

type DataNode struct {
	Key      common.Compareable
	Value    interface{}
	Forwards []*DataNode
}

func NewDataNode(key common.Compareable, value interface{}, level int) *DataNode {
	forwards := make([]*DataNode, level, level)
	return &DataNode{Key: key, Value: value, Forwards: forwards}
}

func (d *DataNode) String() string {
	return fmt.Sprintf("{key: %v, value: %v, level:%d}", d.Key, d.Value, len(d.Forwards))
}

type SkipList struct {
	level       int
	header      *DataNode
	length      uint32
	replaceable bool
}

const (
	MAX_LEVEL = 16
)

type Max int

func (m Max) CompareTo(to interface{}) int {
	return 1
}

var (
	MAX = common.Max(-1)
	NIL = NewDataNode(MAX, "TAIL", 0)
	MIN = common.Min(-1)

	ERR_OUTOFLEVE   = errors.New(fmt.Sprintf("Out of max level == %v", MAX_LEVEL))
	ERR_COMPAREABLE = errors.New("Key should be compareable")

	ERR_EXIST   = errors.New("Key alerady existed")
	ERR_DELETED = errors.New("Key alerady deleted")
)

func NewSkipList(level int, replaceable bool) (*SkipList, error) {
	if level < 0 || level > MAX_LEVEL {
		return nil, ERR_OUTOFLEVE
	}
	header := NewDataNode(MIN, "HEAD", level)
	for i := 0; i < level; i++ {
		header.Forwards[i] = NIL
	}
	sl := &SkipList{level: level, header: header, replaceable: replaceable}
	return sl, nil
}

func (sl *SkipList) Find(key common.Compareable) *DataNode {
	node := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for node.Forwards[i].Key.CompareTo(key) < 0 {
			node = node.Forwards[i]
		}
		if node.Forwards[i].Key.CompareTo(key) == 0 {
			return node.Forwards[i]
		}
	}

	return nil
}

func (sl *SkipList) Insert(key common.Compareable, value interface{}) (*DataNode, error) {
	level := sl.level
	updates := make([]*DataNode, level, level)
	node := sl.header
	for i := level - 1; i >= 0; i-- {
		for node.Forwards[i].Key.CompareTo(key) < 0 {
			node = node.Forwards[i]
		}
		updates[i] = node
	}
	if node.Forwards[0].Key.CompareTo(key) == 0 {
		if sl.replaceable {
			node.Forwards[0].Value = value
			return node.Forwards[0], nil
		} else {
			return nil, ERR_EXIST
		}
	}

	level = sl.randomLevel()
	newNode := NewDataNode(key, value, level)
	for i := level - 1; i >= 0; i-- {
		newNode.Forwards[i] = updates[i].Forwards[i]
		updates[i].Forwards[i] = newNode
	}
	for !atomic.CompareAndSwapUint32(&sl.length, sl.length, sl.length+1) {
	}
	return newNode, nil
}

func (sl *SkipList) Remove(key common.Compareable) *DataNode {
	level := sl.level
	front := make([]*DataNode, level, level)
	node := sl.header
	for i := level - 1; i >= 0; i-- {
		for {
			if node.Forwards[i] == NIL {
				break
			}
			if v := node.Forwards[i].Key.CompareTo(key); v < 0 {
				node = node.Forwards[i]
			} else if v == 0 {
				front[i] = node
				break
			} else {
				break
			}
		}
	}
	if node.Forwards[0].Key.CompareTo(key) != 0 {
		return nil
	}
	node = node.Forwards[0]
	for i := level - 1; i >= 0; i-- {
		if front[i] == nil {
			continue
		}
		front[i].Forwards[i] = node.Forwards[i]
		node.Forwards[i] = nil
	}
	return node
}

func (sl *SkipList) Update(key common.Compareable, value interface{}) *DataNode {
	node := sl.Find(key)
	if node == nil {
		return nil
	}
	node.Value = value
	return node
}

func (sl *SkipList) DisplayAll() {
	node := sl.header
	for {
		node = node.Forwards[0]
		if node != NIL {
			fmt.Println(node)
		} else {
			break
		}
	}
}

func (sl *SkipList) Length() uint32 {
	return sl.length
}

func (sl *SkipList) randomLevel() int {
	num := common.RandInt63()
	num = num ^ (num + 1)
	level := common.NumsOfOne(num)
	if level > sl.level {
		level = sl.level
	}
	return level
}
