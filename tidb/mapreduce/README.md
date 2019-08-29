## Introduction

This is the Map-Reduce homework for PingCAP Talent Plan Online of week 2.

There is a uncompleted Map-Reduce framework, you should complete it and use it to extract the 10 most frequent URLs from data files.

## Getting familiar with the source

The simple Map-Reduce framework is defined in `mapreduce.go`.

It is uncompleted and you should fill your code below comments `YOUR CODE HERE`.

The map and reduce function are defined as same as MIT 6.824 lab 1.
```
type ReduceF func(key string, values []string) string
type MapF func(filename string, contents string) []KeyValue
```

There is an example in `urltop10_example.go` which is used to extract the 10 most frequent URLs.

After completing the framework, you can run this example by `make test_example`.

And then please implement your own `MapF` and `ReduceF` in `urltop10.go` to accomplish this task.

After filling your code, please use `make test_homework` to test.

All data files will be generated at runtime, and you can use `make cleanup` to clean all test data.

Please output URLs by lexicographical order and ensure that your result has the same format as test data so that you can pass all tests.

Each test cases has **different data distribution** and you should take it into account.

## Requirements and rating principles

* (40%) Performs better than `urltop10_example`.
* (20%) Pass all test cases.
* (20%) Profile your program with `pprof`, analyze the performance bottleneck (both the framework and your own code).
* (10%) Have a good code style.
* (10%) Document your idea and code.

NOTE: **go 1.12 is required**

## How to use

Fill your code below comments `YOUR CODE HERE` in `mapreduce.go` to complete this framework.

Implement your own `MapF` and `ReduceF` in `urltop10.go` and use `make test_homework` to test it.

There is a builtin unit test defined in `urltop10_test.go`, however, you still can write your own unit tests.

How to run example:
```
make test_example
```

How to test your implementation:
```
make test_homework
```

How to clean up all test data:
```
make cleanup
```

How to generate test data again:
```
make gendata
```

## 实验笔记

### 准备阶段：阅读框架代码
+ 一个任务结构体：
```go
type task struct {
	dataDir    string
	jobName    string
	mapFile    string   // only for map, the input file
	phase      jobPhase // are we in mapPhase or reducePhase?
	taskNumber int      // this task's index in the current phase
	nMap       int      // number of map tasks
	nReduce    int      // number of reduce tasks
	mapF       MapF     // map function used in this job
	reduceF    ReduceF  // reduce function used in this job
	wg         sync.WaitGroup
}
```
可以看出有目录名字等等，说明对每个文件的操作为一个任务，任务分为`map`任务和`reduce`任务，

+ 一个mapreduce集群结构体：
```go
// MRCluster represents a map-reduce cluster.
type MRCluster struct {
	nWorkers int
	wg       sync.WaitGroup
	taskCh   chan *task
	exit     chan struct{}
}

var singleton = &MRCluster{
	nWorkers: runtime.NumCPU(),
	taskCh:   make(chan *task),
	exit:     make(chan struct{}),
}

func init() {
	singleton.Start()
}
```
`MRCluster`结构体会保存`worker`数量，还有传送`task`的通道，
可以看出该结构体只有一个实例，并且调用该包的时候就开始工作
+ Start函数：
```go
// Start starts this cluster.
func (c *MRCluster) Start() {
	for i := 0; i < c.nWorkers; i++ {
		c.wg.Add(1)
		go c.worker()
	}
}
```
可以看出`Start`函数只是简单地调用`worker`函数
+ worker函数：
```go
func (c *MRCluster) worker() {
	defer c.wg.Done()
	for {
		select {
		case t := <-c.taskCh:
			if t.phase == mapPhase {
                // 省略...
			} else {
				// YOUR CODE HERE :)
				// hint: don't encode results returned by ReduceF, and just output
				// them into the destination file directly so that users can get
				// results formatted as what they want.
				panic("YOUR CODE HERE")
			}
			t.wg.Done()
		case <-c.exit:
			return
		}
	}
}
```
可以看出`worker`函数就是接受从任务通道传来的任务，根据是`map`任务还是`reduce`任务来做不一样的操作
+ Submit函数，这是外部使用的接口
```go
// Submit submits a job to this cluster.
func (c *MRCluster) Submit(jobName, dataDir string, mapF MapF, reduceF ReduceF, mapFiles []string, nReduce int) <-chan []string {
	notify := make(chan []string)
	go c.run(jobName, dataDir, mapF, reduceF, mapFiles, nReduce, notify)
	return notify
}
```
可以看出外部调用需要传一些必要的参数，返回一个通道，用于通知调用者
+ run函数
```go
func (c *MRCluster) run(jobName, dataDir string, mapF MapF, reduceF ReduceF, mapFiles []string, nReduce int, notify chan<- []string) {
	// map phase
	nMap := len(mapFiles)
	tasks := make([]*task, 0, nMap)
	for i := 0; i < nMap; i++ {
		t := &task{
			dataDir:    dataDir,
			jobName:    jobName,
			mapFile:    mapFiles[i],
			phase:      mapPhase,
			taskNumber: i,
			nReduce:    nReduce,
			nMap:       nMap,
			mapF:       mapF,
		}
		t.wg.Add(1)
		tasks = append(tasks, t)
		go func() { c.taskCh <- t }()
	}
	for _, t := range tasks {
		t.wg.Wait()
	}

	// reduce phase
	// YOUR CODE HERE :D
	panic("YOUR CODE HERE")
}
```
可以看出`run`函数就是通过通道将任务分配给`worker`


### 完成第一部分：mapreduce
在`mapreduce.go`中，我们需要补充完reduce部分的代码，包括初始化和执行

先看初始化reduce任务，在`run`函数中：
```go
	// reduce phase
	// YOUR CODE HERE :D
	tasks = make([]*task, 0, nReduce)
	for i := 0; i < nReduce; i++ {
		t := &task{
			dataDir:    dataDir,
			jobName:    jobName,
			phase:      reducePhase,
			taskNumber: i,
			nReduce:    nReduce,
			nMap:       nMap,
			reduceF:    reduceF,
		}
		t.wg.Add(1)
		tasks = append(tasks, t)
		go func() { c.taskCh <- t }()
	}
	for _, t := range tasks {
		t.wg.Wait()
	}

	var inputFiles []string
	for _, t := range tasks {
		inputFiles = append(inputFiles, mergeName(t.dataDir, t.jobName, t.taskNumber))
	}
	notify <- inputFiles
```
参考map初始化的做法，reduce的初始化也就很简单，不过最后不要忘记了要将最后所得到的结果文件要通过通道通知客户端

接着完成真正执行reduce的代码：
```go
				// YOUR CODE HERE :)
				// hint: don't encode results returned by ReduceF, and just output
				// them into the destination file directly so that users can get
				// results formatted as what they want.
				kv_map := make(map[string]([]string))

				for i := 0; i < t.nMap; i++ {
					// 读取map任务产生的文件
					rpath := reduceName(t.dataDir, t.jobName, i, t.taskNumber)
					f, err := os.Open(rpath)
					if err != nil {
						log.Fatalln(err)
					}
					defer f.Close()

					decoder := json.NewDecoder(f)
					var kv KeyValue
					for ; decoder.More(); {
						err := decoder.Decode(&kv)
						if err != nil {
							log.Fatalln(err)
						}
						kv_map[kv.Key] = append(kv_map[kv.Key], kv.Value)
					}
				}

				//keys := make([]string, 0, len(kv_map))
				//for k, _ := range kv_map {
				//	keys = append(keys, k)
				//}
				//sort.Strings(keys)

				// 创建输出文件
				outf, err := os.Create(mergeName(t.dataDir, t.jobName, t.taskNumber))
				if err != nil {
					log.Fatalln(err)
				}
				defer outf.Close()

				for k, v := range kv_map {
					// 写入文件
					outf.WriteString(t.reduceF(k, v))
```
首先要读取由map任务生成的中间文件，并用`kv_map`保存数据，接着创建输出文件，
将由`reduceF`产生的内容写入输出文件

这时候就完成了第一部分的代码，可以执行`make test_example`来测试，
这里有一点要注意，就是由于数据量会很大，特别消耗资源，cpu会占用率达到百分之九十多，
电脑就卡住了，所以最好用虚拟机或云服务器来跑

### 第二部分准备：阅读urltop10_example.go
+ RoundArgs结构体：
```go
// RoundArgs contains arguments used in a map-reduce round.
type RoundArgs struct {
	MapFunc    MapF
	ReduceFunc ReduceF
	NReduce    int
}

// RoundsArgs represents arguments used in multiple map-reduce rounds.
type RoundsArgs []RoundArgs
```
可以看到这个结构体保存的是每一轮mapreduce对应的函数

+ 接着看看例子怎么使用mapreduce的
```go
// ExampleURLTop10 generates RoundsArgs for getting the 10 most frequent URLs.
// There are two rounds in this approach.
// The first round will do url count.
// The second will sort results generated in the first round and
// get the 10 most frequent URLs.
func ExampleURLTop10(nWorkers int) RoundsArgs {
	var args RoundsArgs
	// round 1: do url count
	args = append(args, RoundArgs{
		MapFunc:    ExampleURLCountMap,
		ReduceFunc: ExampleURLCountReduce,
		NReduce:    nWorkers,
	})
	// round 2: sort and get the 10 most frequent URLs
	args = append(args, RoundArgs{
		MapFunc:    ExampleURLTop10Map,
		ReduceFunc: ExampleURLTop10Reduce,
		NReduce:    1,
	})
	return args
}
```
可以看到要计算出现最多的10个url，需要两轮mapreduce，
第一轮map将同一个url弄到一个文件中，第一轮reduce将map生成的中间文件中的每个url汇总成一个url和对应出现次数的字符串；
第二轮map将字符串整合到`KeyValue`结构体的`Value`中，
第二轮reduce先将url和对应出现的次数分开，存到一个哈希表中，然后调用`TopN`排序得到出现次数最多的10个url

### 实现第二部分
在参考了用例代码后，我们要自己实现map和reduce的函数

在之前分析用例代码，我们可以发现，第一轮的map和reduce，还有第二轮的map可以合成一个map任务，
```go
func URLMap(filename string, contents string) []KeyValue {
	lines := strings.Split(string(contents), "\n")
	cnts := make(map[string]int)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		cnts[l]++
	}
	kvs := make([]KeyValue, 0, len(cnts))
	for k, v := range cnts {
		s := fmt.Sprintf("%s %s", k, strconv.Itoa(v))
		kvs = append(kvs, KeyValue{Key: "", Value: s})
	}
	return kvs
}
```

而reduce任务只需要在用例第二轮reduce的基础上改一个地方，就是`=n`改为`+=n`，
因为map会产生多个map文件，而在map任务中，我们的哈希表是在内存中保存数据的，
多次调用这个map函数，哈希表的数据只是限于一个文件的url，map任务生成的文件仍然会有重复的url，
所以要在reduce任务中改为累加
```go
func URLReduce(key string, values []string) string {
	cnts := make(map[string]int)
	for _, v := range values {
		v := strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		tmp := strings.Split(v, " ")
		n, err := strconv.Atoi(tmp[1])
		if err != nil {
			panic(err)
		}
		cnts[tmp[0]] += n
	}

	us, cs := TopN(cnts, 10)
	buf := new(bytes.Buffer)
	for i := range us {
		fmt.Fprintf(buf, "%s: %d\n", us[i], cs[i])
	}
	return buf.String()
}
```

用`make test_homework`测试，效果比用例代码要好一些