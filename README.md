[![Build Status](https://travis-ci.org/mmirolim/gokvc.svg)](https://travis-ci.org/mmirolim/gokvc)
[![GoDoc](https://godoc.org/github.com/mmirolim/gokvc?status.svg)](http://godoc.org/github.com/mmirolim/gokvc)
[![Coverage](http://gocover.io/_badge/github.com/mmirolim/gokvc/cache)](http://gocover.io/github.com/mmirolim/gokvc/cache)
[![Coverage](http://gocover.io/_badge/github.com/mmirolim/gokvc/api)](http://gocover.io/github.com/mmirolim/gokvc/api)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmirolim/gokvc)](https://goreportcard.com/badge/github.com/mmirolim/gokvc)

# gokvc
Fast key-value cache with support of strings, lists and dictionaries for Go with simple HTTP API.

# http api
All params set by using get params. All request has No Cache control headers.
Performance wise [fasthttp](https://github.com/valyala/fasthttp) http library used.
I tried to keep allocation low to reduce gc pressure.
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
