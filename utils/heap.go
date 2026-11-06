package utils

import "container/heap"

// IAnyValue arbitrary value interface
type IAnyValue interface {
	// GetScore retrieves the comparison score
	GetScore() float64
	// SortByAsc specifies the sort order
	//
	// true for ascending order (min-heap)
	//
	// false for descending order (max-heap)
	SortByAsc() bool
}

// Heap for arbitrary data types
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
	// Ascending order, i.e. a min-heap
	if (*avh)[i].SortByAsc() {
		return (*avh)[i].GetScore() < (*avh)[j].GetScore()
	}
	// Descending order, i.e. a max-heap
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
	// Push adds data to the heap
	Push(value IAnyValue)
	// Pop removes data from the heap
	Pop() IAnyValue
}

type HeapStd struct {
	values *anyValueHeap
}

// NewHeap initializes a fresh heap
//
// It can sort slices of arbitrary data types
//
// Implement IAnyValue.SortByAsc to specify the sort order (ascending or descending)
//
// IAnyValue.GetScore provides the numeric basis for sorting
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
