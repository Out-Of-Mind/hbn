package main

import (
		"flag"
		"github.com/out-of-mind/hbn/src/dos"
)

var (
		host string
		workers int
		duration int
)

func init() {
		flag.StringVar(&host, "url", "https://google.com", "usage: -url https://google.com")
		flag.IntVar(&workers, "c", 5, "usage: -c 50 to set 50 threads working")
		flag.IntVar(&duration, "d", 20, "usage: -d 20s to set duration to 20 seconds")
}

func main() {
		flag.Parse()
		d, _ := dos.New(host, "GET", duration, workers, nil, nil)
		d.Run()
}
