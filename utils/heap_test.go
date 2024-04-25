package utils

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Student struct {
	ID    int64
	Name  string
	Score float64
}

func (u Student) SortByAsc() bool {
	return false
}

func (u Student) GetScore() float64 {
	return u.Score
}

func TestAscHeap(t *testing.T) {
	var hh = NewHeap()
	for i := 0; i < 10; i++ {
		student := Student{
			ID:    int64(i),
			Name:  fmt.Sprintf("student_%d", i),
			Score: 1 + 100*rand.Float64(),
		}
		hh.Push(student)
	}

	var current float64
	for i := 0; i < 11; i++ {
		value, ok := hh.Pop().(Student)
		if i == 10 {
			assert.Equal(t, false, ok)
		} else {
			if i == 0 {
				current = value.GetScore()
				continue
			}
			t.Log("score: ", value.GetScore())
			assert.Equal(t, true, current > value.GetScore())
			current = value.GetScore()
		}
	}
}

type Student2 struct {
	ID    int64
	Name  string
	Score float64
}

func (u Student2) SortByAsc() bool {
	return true
}

func (u Student2) GetScore() float64 {
	return u.Score
}

func TestDescHeap(t *testing.T) {
	var hh = NewHeap()
	for i := 0; i < 10; i++ {
		student := Student2{
			ID:    int64(i),
			Name:  fmt.Sprintf("student_%d", i),
			Score: 1 + 100*rand.Float64(),
		}
		hh.Push(student)
	}

	var current float64
	for i := 0; i < 11; i++ {
		value, ok := hh.Pop().(Student2)
		if i == 10 {
			assert.Equal(t, false, ok)
		} else {
			t.Log("score: ", value.GetScore())
			assert.Equal(t, true, current < value.GetScore())
			current = value.GetScore()
		}
	}
}

type Float float64

func (f Float) SortByAsc() bool {
	return true
}

func (f Float) GetScore() float64 {
	return float64(f)
}

func TestAscHeapFloat(t *testing.T) {
	var hh = NewHeap()
	for i := 0; i < 10; i++ {
		hh.Push(Float(i + 1))
	}

	var current float64
	for i := 0; i < 11; i++ {
		value, ok := hh.Pop().(Float)
		if i == 10 {
			assert.Equal(t, false, ok)
		} else {
			t.Log("score: ", value.GetScore())
			assert.Equal(t, true, current < value.GetScore())
			current = value.GetScore()
		}
	}
}

type Int int64

func (i Int) SortByAsc() bool {
	return true
}

func (i Int) GetScore() float64 {
	return float64(i)
}

func TestAscHeapInt(t *testing.T) {
	var hh = NewHeap()
	for i := 0; i < 10; i++ {
		hh.Push(Int(i + 1))
	}

	var current float64
	for i := 0; i < 11; i++ {
		value, ok := hh.Pop().(Int)
		if i == 10 {
			assert.Equal(t, false, ok)
		} else {
			t.Log("score: ", value.GetScore())
			assert.Equal(t, true, current < value.GetScore())
			current = value.GetScore()
		}
	}
}
