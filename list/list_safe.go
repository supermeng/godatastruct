package list

import (
	"fmt"
	// "time"
	"unsafe"

	"sync/atomic"

	"github.com/supermeng/godatastruct/common"
)

type SafeList List

func NewSafeList() *SafeList {
	header := NewNode(MIN, nil)
	header.next = NIL
	return &SafeList{Header: header}
}

func (sl *SafeList) Insert(key common.Compareable, value interface{}) *Node {
	node := sl.Header
	var next *Node
	for {
		next = node.next
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(next))),
			nil, unsafe.Pointer(node.next)) { // bolck
			common.RandomWait()
		}
		if v := next.Key.CompareTo(key); v > 0 { // node =>..(newNode)..=> next
			newNode := NewNode(key, value)
			var forward *Node
			for {
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(node.next))),
					unsafe.Pointer(next), unsafe.Pointer(newNode)) {
					newNode.next = next
					for !atomic.CompareAndSwapUint32(&sl.Length, sl.Length, sl.Length+1) {
					}
					return newNode
				}
				forward = node.next
				for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(forward))),
					nil, unsafe.Pointer(node.next)) {
					common.RandomWait()
				}
				if v := forward.Key.CompareTo(key); v < 0 {
					node = forward
				} else if v == 0 {
					node.Value = value
					return node
				} else {
					next = forward
				}
			}
		} else if v == 0 {
			next.Value = value
			return next
		}
		node = next
	}
}

func (sl *SafeList) Find(key common.Compareable) *Node {
	curr := sl.Header.next
	var forward *Node
	for {
		if v := curr.Key.CompareTo(key); v == 0 {
			return curr
		} else if v > 0 {
			return nil
		}
		forward = curr.next
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(forward))),
			nil, unsafe.Pointer(curr.next)) { // bolck
			common.RandomWait()
		}
		curr = forward
	}
}

func (sl *SafeList) Delete(key common.Compareable) *Node {
	node := sl.Header
	next := node.next
	var forward *Node
	for {
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(next))),
			nil, unsafe.Pointer(node.next)) { // bolck
			common.RandomWait()
		}
		if v := next.Key.CompareTo(key); v == 0 { // node =>...=> next =>...=> forward
			forward = next.next
			for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(forward))),
				nil, unsafe.Pointer(next.next)) {
				common.RandomWait()
			}
			for {
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(node.next))),
					unsafe.Pointer(next), nil) { //unsafe.Pointer(next.next)
					for {
						if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(next.next))),
							unsafe.Pointer(forward), unsafe.Pointer(node)) { // should be ref to nil for delete ?
							node.next = forward
							for !atomic.CompareAndSwapUint32(&sl.Length, sl.Length, sl.Length-1) {
							}
							return next
						}
						forward = next.next
						for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&forward)),
							nil, unsafe.Pointer(next.next)) {
							common.RandomWait()
						}
					}
				}
				forward = node.next
				for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(forward))),
					nil, unsafe.Pointer(node.next)) {
					common.RandomWait()
				}
				if v := forward.Key.CompareTo(key); v < 0 {
					node = forward
				} else if v > 0 {
					return nil
				} else {
					if forward != next { //deleted
						return nil
					}
				}
			}
		} else if v > 0 {
			return nil
		}
		node = next
		next = node.next
	}
}

func (sl *SafeList) Display() {
	node := sl.Header.next
	for {
		if node == NIL {
			fmt.Println()
			return
		}
		node = node.next
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(node))),
			nil, unsafe.Pointer(node.next)) { // bolck
		}
		fmt.Print(node, "    ")
	}
}

func (sl *SafeList) circled() bool {
	slow := sl.Header
	quick := slow.next
	var forward *Node
	for {
		if quick == NIL {
			return false
		}
		forward = slow.next
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(forward))),
			nil, unsafe.Pointer(slow.next)) { // bolck
		}
		slow = forward

		forward = quick.next
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(forward))),
			nil, unsafe.Pointer(quick.next)) { // bolck
		}
		quick = forward
		if quick == NIL {
			return false
		}
		forward = quick.next
		for atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(forward))),
			nil, unsafe.Pointer(quick.next)) { // bolck
		}
		quick = forward
		if quick == slow {
			return true
		}
	}
}

func (sl *SafeList) VerifySafeList() {
	node := sl.Header.next
	for {
		if node == NIL {
			fmt.Println("Verify PASS:", sl.Length)
			return
		}
		if node.next.Key.CompareTo(node.Key) < 0 {
			// fmt.Println(node.next.Key, ":", node.Key)
			sl.Display()
			panic("error")
		}
		node = node.next
	}
}
