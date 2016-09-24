package list

import (
	"fmt"
	"testing"

	"math/rand"
)

type key int

func (k key) CompareTo(to interface{}) int {
	return int(k - to.(key))
}

const (
	THREAD_NUM = 10000

	TEST_NUM = 10

	TEST_TIMES = 1000
)

func singleInsert(sl *SafeList, block chan struct{}) {
	for i := 0; i < TEST_NUM; i++ {
		k := rand.Intn(1000)
		sl.Insert(key(k), rand.Intn(100000))
	}
	block <- struct{}{}
}

func singleDelete(sl *SafeList, block chan struct{}) {
	for i := 0; i < TEST_NUM; i++ {
		k := rand.Intn(1000)
		sl.Delete(key(k))
	}
	block <- struct{}{}
}

func SingleTest() {
	sl := NewSafeList()
	block := make(chan struct{}, THREAD_NUM)
	for i := 0; i < THREAD_NUM; i++ {
		go singleInsert(sl, block)
		go singleDelete(sl, block)
	}
	for i := 0; i < THREAD_NUM; i++ {
		<-block
	}
	// sl.VerifySafeList()
	// for i := 0; i < THREAD_NUM; i++ {
	// 	go singleDelete(sl, block)
	// }
	for i := 0; i < THREAD_NUM; i++ {
		<-block
	}
	// sl.Display()
	sl.VerifySafeList()
}

func Test_SafeList(t *testing.T) {
	for i := 0; i < TEST_TIMES; i++ {
		fmt.Println("test:", i)
		SingleTest()
	}
}
