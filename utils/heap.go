package utils

import "container/heap"

// IAnyValue 任意值接口
type IAnyValue interface {
	// GetScore 获取对比分值
	GetScore() float64
	// SortByAsc 指定排序规则
	//
	// true 为升序(小顶堆)
	//
	// false 为降序(大顶堆)
	SortByAsc() bool
}

// 任意数据类型的堆
type anyValueHeap []IAnyValue

func (avh *anyValueHeap) Push(value any) {
	*avh = append(*avh, value.(IAnyValue))
}

func (avh *anyValueHeap) Pop() any {
	if (*avh).Len() < 1 {
		return nil
	}

	last := (*avh)[(*avh).Len()-1]
	*avh = (*avh)[:(*avh).Len()-1]

	return last
}

func (avh *anyValueHeap) Len() int {
	return len(*avh)
}

func (avh *anyValueHeap) Less(i, j int) bool {
	// 按升序排列，即为小顶堆
	if (*avh)[i].SortByAsc() {
		return (*avh)[i].GetScore() < (*avh)[j].GetScore()
	}
	// 按降序排列，即为大顶堆
	return (*avh)[i].GetScore() > (*avh)[j].GetScore()
}

func (avh *anyValueHeap) Swap(i, j int) {
	if i >= 0 && i < (*avh).Len() && j >= 0 && j < (*avh).Len() {
		(*avh)[i], (*avh)[j] = (*avh)[j], (*avh)[i]
	}
}

func (avh *anyValueHeap) Top() IAnyValue {
	if (*avh).Len() < 1 {
		return nil
	}
	return (*avh)[0]
}

type HeapHelper interface {
	// Push 将数据添加到堆中
	Push(value IAnyValue)
	// Pop 将数据从堆中取出
	Pop() IAnyValue
}

type HeapStd struct {
	values *anyValueHeap
}

// NewHeap 初始化一个全新的堆
//
// 它的作用是可以对任意数据类型的切片进行排序
//
// 实现 IAnyValue.SortByAsc 方法可以指定排序规则（升序、降序）
//
// IAnyValue.GetScore 方法是排序的数值依据
func NewHeap() HeapHelper {
	var values = &anyValueHeap{}
	heap.Init(values)

	return &HeapStd{values: values}
}

func (hs *HeapStd) Push(value IAnyValue) {
	heap.Push(hs.values, value)
}

func (hs *HeapStd) Pop() IAnyValue {
	if hs.values.Len() < 1 {
		return nil
	}
	return heap.Pop(hs.values).(IAnyValue)
}
