package skiplist

import (
	"fmt"
	// "math"
	"math/rand"
	"testing"
)

const (
	test_num = 1600
)

type key int

func (k key) CompareTo(to interface{}) int {
	return int(k - to.(key))
}

func singel_op(sl *SkipList, num int, block chan<- struct{}) {
	size := 10
	for i := 0; i < size; i++ {
		k := rand.Int()
		k = i + num*size
		if _, err := sl.Insert(key(k), k); err != nil {
			fmt.Println("err", err)
			break
		}
	}
	// node := sl.Find(key(5))
	// fmt.Println(node)
	// sl.DisplayAll()
	// node = sl.Remove(key(6))
	// node = sl.Find(key(6))
	// fmt.Println(node)
	// sl.DisplayAll()
	// sl.Update(key(5), 1000)
	// fmt.Println()
	// fmt.Println()

	block <- struct{}{}
}
func Test_SkipList(t *testing.T) {
	fmt.Println(key(36).CompareTo(key(35)))
	// level := 16
	// sl, _ := NewSkipList(level, true)
	// block := make(chan struct{}, 16)

	// for i := 0; i < test_num; i++ {
	// 	go singel_op(sl, i, block)
	// }
	// for i := 0; i < test_num; i++ {
	// 	<-block
	// }
	// sl.DisplayAll()
	// fmt.Println(sl.Length())
}
