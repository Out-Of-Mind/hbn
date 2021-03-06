package hbn

import (
		"os"
		"fmt"
		"time"
		"net/http"
		"io/ioutil"
		"math/rand"
		"sync/atomic"
		"github.com/cheggaaa/pb/v3"
		errs "github.com/out-of-mind/hbn/src/internal/errors"
)

// structure of HBN
type HBN struct {
		url string
		method string
		duration int
		workers int

		useragents []string
		headers map[string]string
		cookies []string

		s chan bool

		errors int
		successes int

		client *http.Client

		totalBytesRead uint64
		minLatency float32
		maxLatency float32
		latency []int64
		rps int64

		use_headers bool
		use_useragents bool
		use_cookies bool

		simple_test bool
}

// making new instance of HBN
func New(url string, method string, duration int, workers int, path_to_config string, use_headers bool, use_useragents bool, use_cookies bool, simple_test bool) *HBN{
		// --- read config ---
		config := getConfig(path_to_config)
		// --- read useragents ---
		useragents := getUseragents(use_useragents, config.Useragents)
		// --- declaring headers ---
		headers := getHeaders(use_headers, config.Headers)
		// --- declaring cookies ---
		cookies := getCookies(use_cookies, config.Cookies)
		// --- making channel to manipulate work of worlers ---
		s := make(chan bool)
		// --- declaring totalBytesRead ---
		totalBytesRead := uint64(0)
		// --- declaring list of latencies
		latency := make([]int64, 1)
		return &HBN{
				url: url,
				method: method,
				duration: duration,
				workers: workers,
				headers: headers,
				cookies: cookies,
				useragents: useragents,
				client: &http.Client{},
				s: s,
				errors: 0,
				successes: 0,
				totalBytesRead: totalBytesRead,
				minLatency: 0.1,
				maxLatency: 0.1,
				latency: latency,
				rps: 0,
				use_headers: use_headers,
				use_useragents: use_useragents,
				use_cookies: use_cookies,
		}
}

// function that user start
func (h *HBN) Run() {
		// starting new bar
		bar := pb.StartNew(h.duration)
		// start testing
		h.start()
		// hello message
		fmt.Printf("Running %ds test for %s\n", h.duration, h.url)
		// bar asynchronous function
		go func() {
				for i := 0; i < h.duration; i++ {
						bar.Increment()
						time.Sleep(time.Second)
				}
		}()
		// waiting for n seconds to end testong
		time.Sleep(time.Duration(h.duration)*time.Second)
		// finish bar working
		bar.Finish()
		// stop testing
		h.stop()
		// --- statistics
		fmt.Printf("total bytes read: %db\n", h.totalBytesRead)
		fmt.Printf("total errors: %d, total success requests: %d\n", h.errors, h.successes)
		// ------
		// finding avarage latency
		avrgLatency := findAvrgLatency(h.latency)
		// converting to seconds
		avrg := float32(avrgLatency)/float32(1e9)
		fmt.Printf("avarage latency is %fs\n", avrg)
		// calculating min and max in list of latencies
		min, max := MinMax(h.latency[1:])
		// converting nanoseconds to seconds
		h.minLatency = convert_nanoseconds_to_seconds(min)
		// converting nanoseconds to seconds
		h.maxLatency = convert_nanoseconds_to_seconds(max)
		// printing minimal and maximum latencies
		fmt.Printf("minimal latancy is: %fs\nmaximum latency is: %fs\n", h.minLatency, h.maxLatency)
		// rps counting and printing
		h.rps = int64((h.errors+h.successes)/h.duration)
		fmt.Printf("rps is %d r/s\n", h.rps)
		// tnx function
		fmt.Println("tnx for using my tool)))")
}

// start DOSing
func (h *HBN) start() {
		// pull into stop channel flase statement
		go func(){h.s<-false}()
		// starting workers
		for i := 0; i < h.workers; i++ {
				go func() {
						//runtime.Gosched()
						// generating seed to provide the real random number
						rand.Seed(time.Now().UnixNano())
						// random useragent
						useragent := h.useragents[rand.Intn(len(h.useragents))]
						// starting infinity loop testing
						for {
								select {
								// if stop channel is true - end testing
								case <- h.s:
										return
								default:
										// if stop channel is true - testing
										// start main sttack function
										err, p := h.attack(useragent)
										// checking for error
										if err != nil {
												// checking if error must to be printed
												if p {
														fmt.Println(err)
														os.Exit(1)
												} else {
														// if error mustn't to be printed - jsut add count of erros
														h.errors += 1
												}
										} else {
												// if wasn't any error - add count of success requests
												h.successes += 1
										}
						}
				}
				}()
		}
}

// stop DOSing
func (h *HBN) stop() {
		// pull into stop channel true statement
		go func(){h.s<-true}()
}

// making http request
func (h *HBN) attack(useragent string) (error, bool) {
		// if method is GET - making get request
		if h.method == "GET" {
				req, err := http.NewRequest(h.method, h.url, nil)
				if err != nil {
						// error, mustn't to be printed
						return err, false
				} else {
						if h.use_useragents {
								// set useragent
								req.Header.Set("User-Agent", useragent)
						}
						if h.use_headers {
								// setting headers from config file
								for key, value := range h.headers {
										req.Header.Set(key, value)
								}
						}
						if h.use_cookies {
								for _, v := range h.cookies {
										req.Header.Add("Set-Cookie", v)
								}
						}
						// start time
						start := time.Now().UnixNano()
						// do request
						resp, err := h.client.Do(req)
						// end time
						end := time.Now().UnixNano()
						h.latency = append(h.latency, (end-start))
						//fmt.Println(end.Sub(start).Round(time.Millisecond))
						if err != nil {
								// error mustn't to be printed
								return err, false
						}
						// new bytes read to count total read bytes
						bytesRead, _ := ioutil.ReadAll(resp.Body)
						// adding total read bytes to counter
						atomic.AddUint64(&h.totalBytesRead, uint64(len(bytesRead)))
						// close body
						resp.Body.Close()
						// no error mustn't to be printed
						return nil, false
				}
		} else {
				// if method isn't allowed - return error which must to be printed
				return errs.ErrMethodDoesNotAllowed(h.method), true
		}
}
