package main

import (
		"flag"
		"github.com/out-of-mind/hbn/src/internal/hbn"
)

var (
		host string
		workers int
		duration int
		path_to_useragents , path_to_headers string
)

func init() {
		flag.StringVar(&host, "url", "https://www.google.com/", "usage: -url https://google.com")
		flag.IntVar(&workers, "c", 5, "usage: -c 50 to set 50 threads working")
		flag.IntVar(&duration, "d", 20, "usage: -d 20s to set duration to 20 seconds")
		flag.StringVar(&path_to_useragents, "u", "./configs/useragents.txt", "usage: -u path/to/your/useragents.txt")
		flag.StringVar(&path_to_headers, "h", "./configs/headers.json", "usage: -h path/tp/your/headers.json")
}

func main() {
		flag.Parse()
		h, _ := hbn.New(host, "GET", duration, workers, path_to_useragents, path_to_headers)
		h.Run()
}
