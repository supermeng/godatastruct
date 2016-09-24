package common

import (
	"time"

	"math/rand"
)

const (
	SLEEP_TIME_NASEC_MASK = 63
)

type Comparator interface {
	Compare(c1, c2 interface{}) int
}

type Compareable interface {
	CompareTo(to interface{}) int
}

func RandIntn(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(n)
}

func RandInt63() int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63()
}

func RandInt() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int()
}

func RandomWait() {
	time.Sleep((time.Duration)(time.Now().Nanosecond()&SLEEP_TIME_NASEC_MASK) * time.Nanosecond)
}

func NumsOfOne(n int64) int {
	c := 0
	for {
		if n == 0 {
			break
		}
		n &= (n - 1)
		c++
	}
	return c
}

type Max int

func (m Max) CompareTo(to interface{}) int {
	return 1
}

type Min int

func (m Min) CompareTo(to interface{}) int {
	return -1
}
