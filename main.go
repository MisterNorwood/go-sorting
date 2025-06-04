package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"
)

// BubbleSort implementation
func BubbleSort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

// InsertionSort implementation
func InsertionSort(arr []int) []int {
	n := len(arr)
	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
	return arr
}

// QuickSort implementation
func QuickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	pivot := arr[len(arr)/2]
	var less, equal, greater []int

	for _, v := range arr {
		switch {
		case v < pivot:
			less = append(less, v)
		case v == pivot:
			equal = append(equal, v)
		case v > pivot:
			greater = append(greater, v)
		}
	}

	sortedLess := QuickSort(less)
	sortedGreater := QuickSort(greater)
	return append(append(sortedLess, equal...), sortedGreater...)
}

// merge function to combine two sorted slices into one sorted slice
func merge(left, right []int) []int {
	result := []int{}
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// Append remaining elements
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// parallelMergeSort sorts the array using parallel merge sort with dynamic thread control
func ParallelMergeSort(arr []int, depth, maxDepth int) []int {
	if len(arr) <= 1 {
		return arr
	}

	// Use sequential mergeSort when depth exceeds max depth
	if depth >= maxDepth {
		return mergeSort(arr)
	}

	// Split array into 4 parts
	quarter := len(arr) / 4
	part1 := arr[:quarter]
	part2 := arr[quarter : 2*quarter]
	part3 := arr[2*quarter : 3*quarter]
	part4 := arr[3*quarter:]

	// Goroutines for parallel sorting of 4 parts
	var wg sync.WaitGroup
	wg.Add(4)

	var sorted1, sorted2, sorted3, sorted4 []int

	// Parallel sorting using dynamic threads based on available CPUs
	go func() {
		defer wg.Done()
		sorted1 = ParallelMergeSort(part1, depth+1, maxDepth)
	}()
	go func() {
		defer wg.Done()
		sorted2 = ParallelMergeSort(part2, depth+1, maxDepth)
	}()
	go func() {
		defer wg.Done()
		sorted3 = ParallelMergeSort(part3, depth+1, maxDepth)
	}()
	go func() {
		defer wg.Done()
		sorted4 = ParallelMergeSort(part4, depth+1, maxDepth)
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Merge the sorted parts sequentially
	return merge(merge(sorted1, sorted2), merge(sorted3, sorted4))
}

// mergeSort (sequential) to use as a fallback when depth is too large
func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := mergeSort(arr[:mid])
	right := mergeSort(arr[mid:])
	return merge(left, right)
}

// Benchmark sorting algorithms with averaging
func BenchmarkSort(sortFunc func([]int) []int, arr []int, iterations int) time.Duration {
	totalTime := time.Duration(0)
	for i := 0; i < iterations; i++ {
		copyArr := make([]int, len(arr))
		copy(copyArr, arr)
		start := time.Now()
		sortFunc(copyArr)
		totalTime += time.Since(start)
	}
	return totalTime / time.Duration(iterations)
}

func BenchmarkSortNoReturn(sortFunc func([]int), arr []int, iterations int) time.Duration {
	totalTime := time.Duration(0)
	for i := 0; i < iterations; i++ {
		copyArr := make([]int, len(arr))
		copy(copyArr, arr)
		start := time.Now()
		sortFunc(copyArr)
		totalTime += time.Since(start)
	}
	return totalTime / time.Duration(iterations)
}

func BenchmarkSortThreaded(sortFunc func([]int, int, int) []int, arr []int, iterations int) time.Duration {
	totalTime := time.Duration(0)
	for i := 0; i < iterations; i++ {
		copyArr := make([]int, len(arr))

		numCPU := runtime.NumCPU()
		maxDepth := numCPU * 2
		copy(copyArr, arr)
		start := time.Now()
		sortFunc(copyArr, 0, maxDepth)
		totalTime += time.Since(start)
	}
	return totalTime / time.Duration(iterations)
}

func main() {
	sizes := []int{100, 1000, 10_000, 100_000}
	iterationsBySize := map[int]int{
		100:     1000,
		1000:    500,
		10_000:  100,
		100_000: 100,
	}

	for _, size := range sizes {
		iterations := iterationsBySize[size]
		fmt.Printf("\nSorting benchmark with %d random integers (averaged over %d runs):\n", size, iterations)

		arr := make([]int, size)
		for i := range arr {
			arr[i] = rand.Intn(size)
		}

		fmt.Println("QuickSort:", BenchmarkSort(QuickSort, arr, iterations))
		fmt.Println("Go's built-in Sort:", BenchmarkSortNoReturn(sort.Ints, arr, iterations))
		fmt.Println("Threaded merge sort:", BenchmarkSortThreaded(ParallelMergeSort, arr, iterations))
		fmt.Println("Bubble Sort:", BenchmarkSort(BubbleSort, arr, iterations))
		fmt.Println("Insertion Sort:", BenchmarkSort(InsertionSort, arr, iterations))
	}
}
