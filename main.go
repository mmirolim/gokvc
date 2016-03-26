package main

import (
	"flag"
	"log"

	"github.com/mmirolim/kvc/api"

	"fmt"

	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", ":8081", "TCP address to listen to")

	// BuildVersion set on build
	BuildVersion = ""
)

func main() {
	fmt.Println("hello this is kvc \n", "I am starting\n")

	m := api.New()

	log.Fatal(fasthttp.ListenAndServe(*addr, m))

}
