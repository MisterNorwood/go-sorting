package main

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBubbleSort(t *testing.T) {
	toSort := BubbleSort([]int{5, 7, 5, 9, 2, 4, 2, 3, 8, 6})
	sorted := []int{2, 2, 3, 4, 5, 5, 6, 7, 8, 9}
	assert.Equal(t, sorted, toSort)
}

func TestInsertionSort(t *testing.T) {
	toSort := InsertionSort([]int{5, 7, 5, 9, 2, 4, 2, 3, 8, 6})
	sorted := []int{2, 2, 3, 4, 5, 5, 6, 7, 8, 9}
	assert.Equal(t, sorted, toSort)
}

func TestParallelMereSort(t *testing.T) {
	numCPU := runtime.NumCPU()
	maxDepth := numCPU * 2
	toSort := ParallelMergeSort([]int{5, 7, 5, 9, 2, 4, 2, 3, 8, 6}, 0, maxDepth)
	sorted := []int{2, 2, 3, 4, 5, 5, 6, 7, 8, 9}
	assert.Equal(t, sorted, toSort)
}

func TestQuickSort(t *testing.T) {
	toSort := QuickSort([]int{5, 7, 5, 9, 2, 4, 2, 3, 8, 6})
	sorted := []int{2, 2, 3, 4, 5, 5, 6, 7, 8, 9}
	assert.Equal(t, sorted, toSort)
}
