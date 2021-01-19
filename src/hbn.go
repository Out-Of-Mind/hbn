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
		use_cookies,  use_useragents, use_headers bool
		simple_test bool
)

func init() {
		flag.StringVar(&host, "url", "https://www.google.com/", "usage: -url https://google.com")
		flag.IntVar(&workers, "w", 5, "usage: -w 50 to set 50 workers working")
		flag.IntVar(&duration, "d", 20, "usage: -d 20s to set duration to 20 seconds")
		flag.StringVar(&path_to_config, "c", "./configs/config.json", "usage: -c path/to/your/config.json")
		flag.StringVar(&method, "m", "GET", "usage: -m GET to set GET method")
		flag.BoolVar(&use_headers, "uh", false, "usage: -uh to read headers from config file")
		flag.BoolVar(&use_useragents, "uu", false, "usage: -uu to read useragents from config file")
		flag.BoolVar(&use_cookies, "uc", false, "usage: -uc to read cookies from config file")
		flag.BoolVar(&simple_test, "s", false, "usage: -s to make simple request without any headers or cookies or useragents")
}

func main() {
		flag.Parse()
		h := hbn.New(host, method, duration, workers, path_to_config, use_headers, use_useragents, use_cookies, simple_test)
		h.Run()
}
