package main

import (
	"fmt"
	// "math"
	"math/rand"
	"time"

	"github.com/supermeng/godatastruct/common"
	"github.com/supermeng/godatastruct/list"
	"github.com/supermeng/godatastruct/skiplist"
)

const (
	test_num = 1000
	size     = 5000
)

type key int

func (k key) CompareTo(to interface{}) int {
	if _, ok := to.(common.Max); ok {
		return -1
	}

	if _, ok := to.(common.Min); ok {
		return 1
	}
	return int(k - to.(key))
}

func insert_op(sl *skiplist.SafeSkipList, num int) {
	for i := 0; i < size; i++ {
		// k := rand.Intn(i + num*size + 1)
		k := rand.Intn(300)
		// k = i + num*size
		if _, err := sl.Insert(key(k), i+num*size+1); err != nil {
			fmt.Println("err", err)
			break
		}
	}
}

func remove_op(sl *skiplist.SafeSkipList, num int) {
	for i := 0; i < size; i++ {
		// k := i + num*size
		// k := rand.Intn((i + num*size + 1))
		// k := rand.Intn(30000)
		k := rand.Intn(300)
		sl.Remove(key(k))
	}
}

func singel_remove(sl *skiplist.SafeSkipList, num int, block chan<- struct{}) {
	remove_op(sl, num)
	block <- struct{}{}
}

func singel_insert(sl *skiplist.SafeSkipList, num int, block chan<- struct{}) {
	insert_op(sl, num)
	block <- struct{}{}
}

func sllength(sl *skiplist.SafeSkipList) uint32 {
	node := sl.Header()
	l := uint32(0)
	for {
		node = node.Forwards[0]
		if node == skiplist.NIL {
			break
		} else {
			l++
		}
	}
	fmt.Println("sllength:", l)
	return l
}

func one_SkipList() bool {
	level := 16
	sl, _ := skiplist.NewSafeSkipList(level, true)
	block := make(chan struct{}, 2*test_num)

	for i := 0; i < test_num; i++ {
		go singel_insert(sl, i, block)
		go singel_remove(sl, i, block)
	}
	for i := 0; i < 2*test_num; i++ {
		<-block
	}
	verifyList(sl)
	sllen := sllength(sl)
	fmt.Println("len:", sl.Length())
	if sl.Length() != sllen {
		sl.DisplayAll()
	}
	return (sl.Length() == sllen)
}

func TestSafeSkipList() {
	for i := 0; i < 100; i++ {
		if one_SkipList() {
			fmt.Println("successed:", i)
		} else {
			fmt.Println("failed:", i)
			return
		}
		time.Sleep(1000 * time.Second)
	}
}

func verifyNode(node *skiplist.DataNode) {
	forwards := node.Forwards
	for id := 1; id < len(forwards); id++ {
		if v := forwards[id].Key.CompareTo(forwards[id-1].Key); v < 0 {
			fmt.Println("id:", forwards[id], "  id-1:", forwards[id-1], " cmp:",
				forwards[id].Key.CompareTo(forwards[id-1].Key), " v:", v)
			for _, u := range forwards {
				fmt.Println(u)
			}
			panic("verified error")
		}
	}
}

func verifyList(sl *skiplist.SafeSkipList) {
	node := sl.Header()
	for {
		node = node.Forwards[0]
		if node == skiplist.NIL {
			break
		} else {
			verifyNode(node)
		}
	}

}

func TestSafeList() {
	for i := 0; i < list.TEST_TIMES; i++ {
		fmt.Println("test:", i)
		list.SingleTest()
	}
}

func main() {
	start := time.Now()
	TestSafeSkipList()
	// TestSafeList()
	total := time.Now().Sub(start).Nanoseconds() / 1000 / 1000
	fmt.Println("duration:", total)
	fmt.Println("avg cost:", total/100)
}
