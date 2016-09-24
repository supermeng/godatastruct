package skiplist

import (
	"fmt"
	"unsafe"

	"sync/atomic"

	"github.com/supermeng/godatastruct/common"
)

var (
	count  = uint32(0)
	ncount = uint32(0)
	ecount = uint32(0)
)

type SafeSkipList SkipList

func NewSafeSkipList(level int, replaceable bool) (*SafeSkipList, error) {
	if level < 0 || level > MAX_LEVEL {
		return nil, ERR_OUTOFLEVE
	}
	header := NewDataNode(MIN, "HEAD", level)
	for i := 0; i < level; i++ {
		header.Forwards[i] = NIL
	}
	sl := &SafeSkipList{level: level, header: header, replaceable: replaceable}
	return sl, nil
}

func (sl *SafeSkipList) Find(key common.Compareable) *DataNode {
	level := sl.level
	node := sl.header
	var forward *DataNode
	for i := level - 1; i >= 0; i-- {
		forward = node.Forwards[i]
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&forward)),
			nil, unsafe.Pointer(node.Forwards[i])) {
			common.RandomWait()
		}
		for {
			if v := forward.Key.CompareTo(key); v < 0 {
				node = forward
				forward = node.Forwards[i]
				for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&forward)),
					nil, unsafe.Pointer(node.Forwards[i])) {
					common.RandomWait()
				}
			} else if v == 0 {
				return forward
			} else {
				break
			}
		}
	}
	return nil
}

func (sl *SafeSkipList) Insert(key common.Compareable, value interface{}) (*DataNode, error) {
	level := sl.level
	updates := make([]*DataNode, level, level)
	nexts := make([]*DataNode, level, level)
	node := sl.header
	var forward *DataNode
	for i := level - 1; i >= 0; i-- {
		forward = node.Forwards[i]
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&forward)),
			nil, unsafe.Pointer(node.Forwards[i])) {
			common.RandomWait()
		}
		for {
			if v := forward.Key.CompareTo(key); v < 0 {
				node = forward
				forward = node.Forwards[i]
				for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&forward)),
					nil, unsafe.Pointer(node.Forwards[i])) {
					common.RandomWait()
				}
			} else if v == 0 {
				if sl.replaceable {
					forward.Value = value
					return forward, nil
				}
				return nil, ERR_EXIST
			} else {
				updates[i] = node
				nexts[i] = forward
				break
			}
		}
	}

	if forward.Key.CompareTo(key) == 0 {
		if sl.replaceable {
			forward.Value = value
			return forward, nil
		}
		return nil, ERR_EXIST
	}
	level = sl.randomLevel()
	newNode := NewDataNode(key, value, level)
	for i := level - 1; i >= 0; i-- {
		if node, err := sl.insertLinkIndex(updates[i], nexts[i], newNode, i); err != nil {
			return nil, err
		} else if node != nil {
			// minL := len(node.Forwards)
			// if minL > len(newNode.Forwards) {
			// 	minL = len(newNode.Forwards)
			// }
			// if i < minL-1 {
			// 	fmt.Println(node, ",", newNode, ",", i)
			// 	for _, up := range node.Forwards {
			// 		fmt.Println("up:", up)
			// 	}
			// 	for _, up := range newNode.Forwards {
			// 		fmt.Println("newup:", up)
			// 	}
			// 	panic("what?")
			// }
			for j := level - 1; j > i; j-- {
				if err := sl.deleteLinkIndex(updates[j], nexts[j], newNode, j); err != nil {
					fmt.Println("key:", key, ":", newNode.Forwards[j], " index:", j)
					panic("shit happend!")
					// if sl.replaceable {
					// 	node.Value = value
					// 	return node, nil
					// }
					// return nil, ERR_EXIST
				}
			}
			for j := i; j >= 0; j-- {
				newNode.Forwards[j] = sl.header
			}
			if sl.replaceable {
				node.Value = value
				return node, nil
			}
			return nil, ERR_EXIST
		}
	}
	for !atomic.CompareAndSwapUint32(&sl.length, sl.length, sl.length+1) {
	}
	// fmt.Println("newnode:", newNode)
	return newNode, nil
}

func (sl *SafeSkipList) Remove(key common.Compareable) *DataNode {
	level := sl.level
	updates := make([]*DataNode, level, level)
	nexts := make([]*DataNode, level, level)
	nodes := make([]*DataNode, level, level)
	node := sl.header
	var forward *DataNode
	for i := level - 1; i >= 0; i-- {
		forward = node.Forwards[i]
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&forward)),
			nil, unsafe.Pointer(node.Forwards[i])) {
			common.RandomWait()
		}
		for {
			if v := forward.Key.CompareTo(key); v < 0 {
				node = forward
				forward = node.Forwards[i]
				for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&forward)),
					nil, unsafe.Pointer(node.Forwards[i])) {
					common.RandomWait()
				}
			} else if v == 0 {
				updates[i] = node
				nexts[i] = forward.Forwards[i]
				nodes[i] = forward
				for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&nexts[i])),
					nil, unsafe.Pointer(forward.Forwards[i])) {
					common.RandomWait()
				}
				break
			} else {
				nodes[i] = forward
				break
			}
		}
	}
	if forward.Key.CompareTo(key) != 0 {
		for !atomic.CompareAndSwapUint32(&ncount, ncount, ncount+1) {
		}
		return nil
	}
	node = forward
	level = len(node.Forwards)
	// if updates[level-1] == nil {
	// 	for !atomic.CompareAndSwapUint32(&count, count, count+1) {
	// 	}
	// 	return nil
	// }
	for i := level - 1; i > 0; i-- {
		if nodes[i] != nodes[i-1] {
			for !atomic.CompareAndSwapUint32(&count, count, count+1) {
			}
			return nil
		}
	}
	if updates[level-1] == nil {
		panic("wtf!!")
	}

	for i := level - 1; i >= 0; i-- {
		if updates[i] == nil {
			fmt.Println("key:", key)
			for j := level - 1; j >= 0; j-- {
				fmt.Println("updates:", updates[j], "  next:", nodes[j])
			}
			panic("wao")
			return nil
		}
		if err := sl.deleteLinkIndex(updates[i], nexts[i], node, i); err != nil {
			for !atomic.CompareAndSwapUint32(&ecount, ecount, ecount+1) {
			}
			return nil
		}
	}
	for !atomic.CompareAndSwapUint32(&sl.length, sl.length, sl.length-1) {
	}
	// fmt.Println("remove:", node)
	return node
}

func (sl *SafeSkipList) Update(key common.Compareable, value interface{}) *DataNode {
	node := sl.Find(key)
	if node == nil {
		return nil
	}
	node.Value = value
	return node
}

func (sl *SafeSkipList) Header() *DataNode {
	return sl.header
}

func (sl *SafeSkipList) Length() uint32 {
	return sl.length
}

func (sl *SafeSkipList) randomLevel() int {
	num := common.RandInt63()
	num = num ^ (num + 1)
	level := common.NumsOfOne(num)
	if level > sl.level {
		level = sl.level
	}
	return level
}

func (sl *SafeSkipList) DisplayLevel(index int) {
	node := sl.header
	fmt.Print("level", index, ":  ")
	for {
		node = node.Forwards[index]
		if node == nil {
			panic("hehe")
		}
		if node != NIL {
			fmt.Print(node.Key, "=>")
		} else {
			break
		}
	}
	fmt.Println()
}

func (sl *SafeSkipList) DisplayAll() {
	fmt.Println("ncount:", atomic.LoadUint32(&ncount))
	fmt.Println("ecount:", atomic.LoadUint32(&ecount))
	fmt.Println("count:", atomic.LoadUint32(&count))
	for i := sl.level - 1; i >= 0; i-- {
		sl.DisplayLevel(i)
	}
}

func (sl *SafeSkipList) insertLinkIndex(front, next, newNode *DataNode, index int) (*DataNode, error) {
	i := index
	var forward *DataNode
	key := newNode.Key
	for {
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(front.Forwards[i]))),
			unsafe.Pointer(next), unsafe.Pointer(newNode)) {
			newNode.Forwards[i] = next
			break
		}
		forward = front.Forwards[i]
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&forward)),
			nil, unsafe.Pointer(front.Forwards[i])) {
			common.RandomWait()
		}
		if v := forward.Key.CompareTo(key); v > 0 {
			next = forward
		} else if v == 0 {
			return forward, nil
		} else {
			front = forward
		}
	}
	return nil, nil
}

func (sl *SafeSkipList) deleteLinkIndex(front, next, node *DataNode, index int) error {
	var forward *DataNode
	i := index
	key := node.Key
	for {
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(front.Forwards[i]))),
			unsafe.Pointer(node), nil) {
			for {
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(node.Forwards[i]))),
					unsafe.Pointer(next), unsafe.Pointer(sl.header)) { // should be ref to nil for delete ?
					front.Forwards[i] = next
					return nil
				}
				next = node.Forwards[i]
				for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&next)),
					nil, unsafe.Pointer(node.Forwards[i])) {
					common.RandomWait()
				}
			}
		}
		forward = front.Forwards[i]
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&forward)),
			nil, unsafe.Pointer(front.Forwards[i])) {
			common.RandomWait()
		}
		if v := forward.Key.CompareTo(key); v < 0 {
			front = forward
		} else if v > 0 { // have bean alerady deleted
			return ERR_DELETED
		} else {
			if forward != node { //deleted
				return ERR_DELETED
			}
		}
	}
	return nil
}
