package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// URLTop10 .
func URLTop10(nWorkers int) RoundsArgs {
	// YOUR CODE HERE :)
	// And don't forget to document your idea.
	var args RoundsArgs
	args = append(args, RoundArgs{
		MapFunc:    URLMap,
		ReduceFunc: URLReduce,
		NReduce:    1,
	})
	return args
}

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