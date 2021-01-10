package main

import (
		"flag"
		"github.com/out-of-mind/hbn/src/internal/hbn"
)

var (
		host string
		workers int
		duration int
		path_to_config string
		method string
)

func init() {
		flag.StringVar(&host, "url", "https://www.google.com/", "usage: -url https://google.com")
		flag.IntVar(&workers, "w", 5, "usage: -c 50 to set 50 workers working")
		flag.IntVar(&duration, "d", 20, "usage: -d 20s to set duration to 20 seconds")
		flag.StringVar(&path_to_config, "c", "./configs/config.json", "usage: -w path/to/your/config.json")
		flag.StringVar(&method, "m", "GET", "usage: -m GET to set GET method")
}

func main() {
		flag.Parse()
		h, _ := hbn.New(host, method, duration, workers, path_to_config)
		h.Run()
}
