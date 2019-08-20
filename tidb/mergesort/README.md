## Introduction

This is the Merge Sort home work for PingCAP Talent Plan Online of week 1.

There are 16, 000, 000 int64 values stored in an unordered array. Please
supplement the `MergeSort()` function defined in `mergesort.go` to sort this
array.

Requirements and rating principles:
* (30%) Pass the unit test.
* (20%) Performs better than `sort.Slice()`.
* (30%) Profile your program with `pprof`, analyze the performance bottleneck.
* (10%) Have a good code style.
* (10%) Document your idea and code.

NOTE: **go 1.12 is required**

## How to use

Please supplement the `MergeSort()` function defined in `mergesort.go` to accomplish
the home work.

**NOTE**:
1. There is a builtin unit test defined in `mergesort_test.go`, however, you still
   can write your own unit tests.
2. There is a builtin benchmark test defined in `bench_test.go`, you should run
   this benchmark to ensure that your parallel merge sort is fast enough.


How to test:
```
make test
```

How to benchmark:
```
make bench
```

## 实验笔记
### version1.0：顺序mergesort
首先当然要先实现顺序执行的归并排序，熟悉归并排序的话应该很快就能写出来：
```golang
func MergeSort(src []int64) {
	if len(src) > 1 {
		mid := len(src) / 2
		MergeSort(src[0:mid])
		MergeSort(src[mid:])
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
```

测试一下性能：
```
$ make bench
go test -bench Benchmark -run xx -count 5 -benchmem
goos: linux
goarch: amd64
BenchmarkMergeSort-4                   1        3652349038 ns/op        3221227456 B/op 16777220 allocs/op
BenchmarkMergeSort-4                   1        3633829349 ns/op        3221225472 B/op 16777215 allocs/op
BenchmarkMergeSort-4                   1        3469851788 ns/op        3221225568 B/op 16777216 allocs/op
BenchmarkMergeSort-4                   1        3519544317 ns/op        3221225568 B/op 16777216 allocs/op
BenchmarkMergeSort-4                   1        3509620377 ns/op        3221227360 B/op 16777219 allocs/op
BenchmarkNormalSort-4                  1        4326288097 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4328225099 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4336005161 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4329031691 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4307245016 ns/op              64 B/op          2 allocs/op
PASS
ok      _/home/liu/Desktop/programing/go/src/mergesort  44.410s
```
可以看出我们实现的归并排序时间上比`sort.Slice()`快一点，空间上则要多很多，这也是归并排序的特点，下面我们来使用多goroutine。

### version2.0：并发mergesort
```golang
const max = 1 << 11

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
			wg.Add(2)

			go func() {
				defer wg.Done()
				MergeSort(src[0:mid])
			} ()

			go func() {
				defer wg.Done()
				MergeSort(src[mid:])
			} ()

			wg.Wait()
		}
		merge(src, mid)
	}
}
```
merge函数不变，当切片长度小于max就采用顺序归并排序，否则用两个goroutine来分别处理切片的前半部分和后半部分

测试一下性能：
```
$ make bench
go test -bench Benchmark -run xx -count 5 -benchmem
goos: linux
goarch: amd64
BenchmarkMergeSort-4                   1        1357407428 ns/op        3223403536 B/op 16795017 allocs/op
BenchmarkMergeSort-4                   1        1278981201 ns/op        3221835600 B/op 16791592 allocs/op
BenchmarkMergeSort-4                   1        1207914466 ns/op        3221748464 B/op 16790498 allocs/op
BenchmarkMergeSort-4                   1        1273964833 ns/op        3221769264 B/op 16790794 allocs/op
BenchmarkMergeSort-4                   1        1381991556 ns/op        3222084400 B/op 16793185 allocs/op
BenchmarkNormalSort-4                  1        4251177870 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4327126168 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4335611459 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4298737534 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4267104863 ns/op              64 B/op          2 allocs/op
PASS
ok      _/home/liu/Desktop/programing/go/src/mergesort  33.034s
```
可以看出，时间上有了很大的优化，已经达到了`sort.Slice()`三分之一还少一点，其实到这就差不多了，不过还有一个小点可以再优化一下。

### version3.0：优化
```golang
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

			//直接用主goroutine来调用，减少goroutine的消耗
			MergeSort(src[mid:])

			wg.Wait()
		}
		merge(src, mid)
	}
}
```
改动很少，只是用该函数的主goroutine来直接处理后半部分的mergesort，而不用另起一个goroutine

测试一下性能：
```
$ make bench
go test -bench Benchmark -run xx -count 5 -benchmem
goos: linux
goarch: amd64
BenchmarkMergeSort-4                   1        1483127005 ns/op        3222228304 B/op 16789560 allocs/op
BenchmarkMergeSort-4                   1        1233106729 ns/op        3221494512 B/op 16787346 allocs/op
BenchmarkMergeSort-4                   1        1207362253 ns/op        3221504560 B/op 16786989 allocs/op
BenchmarkMergeSort-4                   1        1188559496 ns/op        3221479280 B/op 16787110 allocs/op
BenchmarkMergeSort-4                   1        1219738624 ns/op        3221482992 B/op 16787024 allocs/op
BenchmarkNormalSort-4                  1        4226798660 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4412368726 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4568287088 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4472948370 ns/op              64 B/op          2 allocs/op
BenchmarkNormalSort-4                  1        4361268809 ns/op              64 B/op          2 allocs/op
PASS
ok      _/home/liu/Desktop/programing/go/src/mergesort  33.561s
```
较上一个版本没有太大的变化，可能是因为减少的goroutine不多吧。。。
