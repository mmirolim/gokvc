package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/mmirolim/gokvc/api"
	"github.com/mmirolim/gokvc/cache"

	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", ":8081", "TCP address to listen to")

	// BuildVersion set on build
	BuildVersion = ""
)

func main() {
	flag.Parse()
	defer glog.Flush()

	glog.Infof("hello this is gokvc \nI am starting\n on port %s", *addr)
	glog.Infof("Build Version %s", BuildVersion)

	// init cache
	cache.Init()

	m := api.New()

	if err := fasthttp.ListenAndServe(*addr, m); err != nil {
		glog.Fatal(err)
	}
}
