package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortingAlgorithms(t *testing.T) {
	type sorterFunc struct {
		name string
		sort func([]int) []int
	}

	sortingAlgorithms := []sorterFunc{
		{"BubbleSort", BubbleSort},
		{"InsertionSort", InsertionSort},
		{"ParallelMergeSort", ParallelMergeSort},
		{"HeapSort", HeapSort},
		{"BucketSort", BucketSort},
		{"QuickSort", QuickSort},
	}

	input := []int{5, 7, 5, 9, 2, 4, 2, 3, 8, 6}
	expected := []int{2, 2, 3, 4, 5, 5, 6, 7, 8, 9}

	for _, algo := range sortingAlgorithms {
		t.Run(algo.name, func(t *testing.T) {
			sorted := algo.sort(append([]int(nil), input...))
			assert.Equal(t, expected, sorted)
		})
	}

	t.Run("RadixSortSigned", func(t *testing.T) {
		in := append([]int(nil), input...)
		RadixSortSigned(in)
		assert.Equal(t, expected, in)
	})
}
