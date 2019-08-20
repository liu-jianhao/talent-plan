package main

import "sync"

const max = 1 << 11

// MergeSort performs the merge sort algorithm.
// Please supplement this function to accomplish the home work.
func MergeSort(src []int64) {
	if len(src) > 1 {
		mid := len(src) / 2
		// 定一个阀值，小于这个阀值就不再用多goroutine来MergeSort
		// 如果不定一个阀值的话goroutine会过多，导致系统卡住，别问我是怎么知道的-_-
		if len(src) <= max {
			MergeSort(src[0:mid])
			MergeSort(src[mid:])
		} else {
			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				MergeSort(src[0:mid])
			} ()

			//go func() {
			//	defer wg.Done()
			//	MergeSort(src[mid:])
			//} ()
			//直接用主goroutine来调用，减少goroutine的消耗
			MergeSort(src[mid:])

			wg.Wait()
		}
		merge(src, mid)
	}
}

func merge(src []int64, mid int) {
	cpy := make([]int64, len(src))
	copy(cpy, src)

	left := 0
	right := mid
	cur := 0
	end := len(src) - 1

	for left <= mid-1 && right <= end {
		if cpy[left] <= cpy[right] {
			src[cur] = cpy[left]
			left++
		} else {
			src[cur] = cpy[right]
			right++
		}
		cur++
	}

	for left <= mid - 1 {
		src[cur] = cpy[left]
		left++
		cur++
	}
}
