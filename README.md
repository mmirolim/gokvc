[![Build Status](https://travis-ci.org/mmirolim/gokvc.svg)](https://travis-ci.org/mmirolim/gokvc)
[![GoDoc](https://godoc.org/github.com/mmirolim/gokvc?status.svg)](http://godoc.org/github.com/mmirolim/gokvc)
[![Coverage](https://gocover.io/_badge/github.com/mmirolim/gokvc/cache)](https://gocover.io/github.com/mmirolim/gokvc/cache)
[![Coverage](https://gocover.io/_badge/github.com/mmirolim/gokvc/api)](https://gocover.io/github.com/mmirolim/gokvc/api)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmirolim/gokvc)](https://goreportcard.com/badge/github.com/mmirolim/gokvc)

# gokvc
Fast key-value cache with support of strings, lists, dictionaries and per-key TTL for Go with convenient HTTP API.
I tried to keep allocation low to reduce gc pressure. All cache buckets use concurrent stripped maps to reduce lock contention. Reflections and interface{} usage avoided as much as possible to improve performance. Expired keys removed by concurrent gc in batches.

## Work in progress edge cases, api changes expected

## TODO
- Native client library
- Increase test coverage
- Improve performance
- Scaling
- Persistence

# Installation
	
	cd $GOPATH/src/github.com/mmirolim
	git clone git@github.com:mmirolim/gokvc.git
	cd gokvc
	make run
	
# Usage
	
	gokvc -addr=":8081" -log_dir="logs" -stderrthreshold=INFO -v=3
	
# Testing

	make test

# http api
All params set by using get params. All request has No Cache control headers.
Performance wise [fasthttp](https://github.com/valyala/fasthttp) http library used.

## Strings

	get the value of a key
	ADDR /get?k={key} // val
	
	set the string value of a key
	ADDR /set?k={key}&v={val}
	
	removes the specified key. A key is ignored if not exist
	ADDR /del?k={key}
	
	get the ttl (in seconds) left of a not expired key
	ADDR /ttl?k={key}
	
	returns all not expired keys
	ADDR /keys // [key1, key2, ...]
	
	returns number of not expired strings
	ADDR /slen
	
## Lists

	get all values of list by a key
	ADDR /lget?k={key} // [val1, val2, ...]
	
	prepend a val to a list. If not exist, list created
	ADDR /lpush?k={key}&v={val}&t={ttl} 
	
	remove and get the first element in a list. If empty removes list
	ADDR /lpop?k={key} // val
	
	removes the specified list. A key is ignored if not exist
	ADDR /ldel?k={key}
	
	get the ttl (in seconds) left of a not expired key
	ADDR /lttl?k={key}
	
	returns all not expired keys of lists
	ADDR /lkeys // [key1, key2, ...]
	
	returns number of not expired lists
	ADDR /lslen
	
## Dictionaries

	get all values of list by a key
	ADDR /dget?k={key} // map[field1: val1, field2: val2, ...]
	
	prepend a val to a list
	ADDR /dfget?k={key}&f={field} // val
	
	set field to dictionary with defined key. If not exist dictionary created
	ADDR /dfset?k={key}&f={field}&v={val}&t={ttl}
	
	remove dictionary. A key is ignored if not exist
	ADDR /ddel?k={key}
	
	remove entry from dictionary
	ADDR /dfdel?k={key}&f={field}
	
	get the ttl (in seconds) left of a not expired key
	ADDR /dttl?k={key}
	
	returns all not expired keys of lists
	ADDR /dkeys // [key1, key2, ...]
	
	returns number of not expired lists
	ADDR /dslen

# Benchmarks
Reuslt with Go 1.6 on CPU Core™ i7-5500U @ 2.40GHz × 4 (Fedora23)
## Http API benchmark

	go test ./cache -run=none -bench=^Benchmark -cpu=1,3
	
	BenchmarkTimeNow      	50000000	        23.7 ns/op	       0 B/op	       0 allocs/op
	BenchmarkTimeNow-3    	50000000	        24.5 ns/op	       0 B/op	       0 allocs/op
	BenchmarkSysTime      	500000000	         3.45 ns/op	       0 B/op	       0 allocs/op
	BenchmarkSysTime-3    	500000000	         3.64 ns/op	       0 B/op	       0 allocs/op
	BenchmarkGetTtl       	20000000	        72.1 ns/op	       0 B/op	       0 allocs/op
	BenchmarkGetTtl-3     	20000000	        73.5 ns/op	       0 B/op	       0 allocs/op
	BenchmarkSet          	20000000	        80.2 ns/op	       0 B/op	       0 allocs/op
	BenchmarkSet-3        	20000000	        82.2 ns/op	       0 B/op	       0 allocs/op
	BenchmarkGet          	20000000	        74.7 ns/op	       0 B/op	       0 allocs/op
	BenchmarkGet-3        	20000000	        76.7 ns/op	       0 B/op	       0 allocs/op
	BenchmarkDel          	20000000	        91.1 ns/op	       0 B/op	       0 allocs/op
	BenchmarkDel-3        	20000000	        93.6 ns/op	       0 B/op	       0 allocs/op
	BenchmarkParallelSET  	20000000	        84.5 ns/op
	BenchmarkParallelSET-3	20000000	        76.0 ns/op
	BenchmarkParallelGET  	20000000	        74.8 ns/op
	BenchmarkParallelGET-3	20000000	        71.2 ns/op
	ok  	github.com/mmirolim/gokvc/cache	26.749s
	
	// github.com/rakyll/boom
	boom -n 10000 -c 50 http://localhost:8081/get?k=key1
	
	Summary:
	Average:	0.0012 secs
	Requests/sec:	41456.5449
	Total data:	315648 bytes
	Size/request:	31 bytes



## Cache API benchmarks 

	go test ./api -run=none -bench=^Benchmark -cpu=1,3
	
	BenchmarkSET  	10000000	       211 ns/op	       0 B/op	       0 allocs/op
	BenchmarkSET-3	10000000	       224 ns/op	       0 B/op	       0 allocs/op
	BenchmarkGET  	10000000	       124 ns/op	       0 B/op	       0 allocs/op
	BenchmarkGET-3	10000000	       129 ns/op	       0 B/op	       0 allocs/op
	BenchmarkTTL  	10000000	       124 ns/op	       0 B/op	       0 allocs/op
	BenchmarkTTL-3	10000000	       129 ns/op	       0 B/op	       0 allocs/op
	ok  	github.com/mmirolim/gokvc/api	10.436s
