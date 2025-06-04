package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

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

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
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
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := mergeSort(arr[:mid])
	right := mergeSort(arr[mid:])
	return merge(left, right)
}

func ParallelMergeSort(arr []int) []int {
	type task struct {
		data  []int
		depth int
	}

	maxDepth := 4

	var sort func([]int, int) []int
	sort = func(arr []int, depth int) []int {
		if len(arr) <= 2048 || depth >= maxDepth {
			return mergeSort(arr)
		}

		mid := len(arr) / 2
		var left, right []int
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			left = sort(arr[:mid], depth+1)
		}()

		right = sort(arr[mid:], depth+1)
		wg.Wait()

		return merge(left, right)
	}

	return sort(arr, 0)
}

func BenchmarkSort(sortFn func([]int) []int, base []int, iterations int) int64 {
	var total int64
	for i := 0; i < iterations; i++ {
		a := make([]int, len(base))
		copy(a, base)

		start := time.Now()
		sortFn(a)
		total += time.Since(start).Microseconds()
	}
	return total / int64(iterations)
}

func BenchmarkSortNoReturn(sortFn func([]int), base []int, iterations int) int64 {
	var total int64
	for i := 0; i < iterations; i++ {
		a := make([]int, len(base))
		copy(a, base)

		start := time.Now()
		sortFn(a)
		total += time.Since(start).Microseconds()
	}
	return total / int64(iterations)
}

func RadixSortSigned(nums []int) {
	if len(nums) == 0 {
		return
	}

	var negs, poss []int
	for _, n := range nums {
		if n < 0 {
			negs = append(negs, -n)
		} else {
			poss = append(poss, n)
		}
	}

	radixSortBase10(negs)
	radixSortBase10(poss)

	for i := 0; i < len(negs)/2; i++ {
		negs[i], negs[len(negs)-1-i] = negs[len(negs)-1-i], negs[i]
	}
	for i := range negs {
		negs[i] = -negs[i]
	}

	copy(nums, append(negs, poss...))
}

func radixSortBase10(nums []int) {
	if len(nums) == 0 {
		return
	}

	max := nums[0]
	for _, n := range nums {
		if n > max {
			max = n
		}
	}

	exp := 1
	for max/exp > 0 {
		count := make([]int, 10)
		output := make([]int, len(nums))

		for _, n := range nums {
			d := (n / exp) % 10
			count[d]++
		}
		for i := 1; i < 10; i++ {
			count[i] += count[i-1]
		}
		for i := len(nums) - 1; i >= 0; i-- {
			d := (nums[i] / exp) % 10
			count[d]--
			output[count[d]] = nums[i]
		}
		copy(nums, output)
		exp *= 10
	}
}

func main() {
	sizes := []int{100, 1000, 10_000, 100_000, 1_000_000}
	iterationsBySize := map[int]int{
		100:       1000,
		1000:      500,
		10_000:    100,
		100_000:   10,
		1_000_000: 5,
	}

	file, err := os.Create("results.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Size", "Algorithm", "Time (ms)"})

	for _, size := range sizes {
		iterations := iterationsBySize[size]
		fmt.Printf("\nSorting benchmark with %d random integers (averaged over %d runs):\n", size, iterations)

		arr := make([]int, size)
		for i := range arr {
			arr[i] = rand.Intn(size)
		}

		bench := func(name string, f func([]int) []int) {
			ms := BenchmarkSort(f, arr, iterations)
			fmt.Printf("%s: %.2f ms\n", name, ms)
			writer.Write([]string{strconv.Itoa(size), name, fmt.Sprintf("%.2f", ms)})
		}
		benchNR := func(name string, f func([]int)) {
			ms := BenchmarkSortNoReturn(f, arr, iterations)
			fmt.Printf("%s: %.2f ms\n", name, ms)
			writer.Write([]string{strconv.Itoa(size), name, fmt.Sprintf("%.2f", ms)})
		}

		bench("QuickSort", QuickSort)
		benchNR("Go's built-in Sort", sort.Ints)
		benchNR("Radix Sort", RadixSortSigned)
		bench("Threaded merge sort", ParallelMergeSort)
		bench("Insertion Sort", InsertionSort)
		bench("Bubble Sort", BubbleSort)
	}
}
