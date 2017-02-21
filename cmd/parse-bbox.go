package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-bbox/parser"
	"log"
)

func main() {

	var bbox = flag.String("bbox", "", "...")
	var scheme = flag.String("scheme", "swne", "...")
	var order = flag.String("order", "latlon", "...")

	flag.Parse()

	p, err := parser.NewParser()

	if err != nil {
		log.Fatal(err)
	}

	p.Scheme = *scheme
	p.Order = *order

	bb, err := p.Parse(*bbox)

	fmt.Print(bb)
}
