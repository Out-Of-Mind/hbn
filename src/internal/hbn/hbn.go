package hbn

import (
		"os"
		"fmt"
		"time"
		"bufio"
		"runtime"
		"strings"
		"net/http"
		"io/ioutil"
		"math/rand"
		"sync/atomic"
		"encoding/json"
		"github.com/cheggaaa/pb/v3"
)

// structure of HBN
type HBN struct {
		url string
		method string
		duration int
		workers int
		useragents []string
		headers map[string]string
		s chan bool
		errors int
		succeses int
		client *http.Client
		totalBytesRead uint64
		minLatency float32
		maxLatency float32
		latency []int64
}

// making new instance of HBN
func New(url string, method string, duration int, workers int, path_to_config string) (*HBN, error) {
		// --- read config ---
		// opening config file
		c, err := os.Open(path_to_config)
		if err != nil {
				// if was error - handling it
				fmt.Println(ErrConfigWasNotFound(path_to_config))
		}
		decoder := json.NewDecoder(c)
		// new instance of Config structure to decoding json
		config := new(Config)
		// decoding json
		err = decoder.Decode(&config)
		if err != nil {
				// if was error - handling it
				fmt.Println(ErrConfigRead())
				os.Exit(1)
		}
		// close file
		c.Close()
		// ------
		// --- read useragents ---
		// opening file where useragents
		f, err := os.Open(config.Useragents)
		if err != nil {
				// if was error - handling it
				fmt.Println(ErrUseragentsWasNotFound(config.Useragents))
				os.Exit(1)
		}
		// making new scanner
		scanner := bufio.NewScanner(f)
		// making list of useragents
		useragents := make([]string, 1)
		for scanner.Scan() {
				// reading useragents line bt line
				useragents = append(useragents, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
				// if was error - handling it
				fmt.Println(ErrReadFile("useragents"))
				os.Exit(1)
		}
		// close file
		f.Close()
		// ------
		// --- declaring headers ---
		// making new map - headera
		headers := make(map[string]string)
		// enumerating headers from config
		for _, v := range config.Headers {
				// split header
				o := strings.Split(v, " ")
				// setting header from splitted string
				headers[o[0]] = o[1]
		}
		// ------
		// --- making channel to manipulate work of worlers ---
		s := make(chan bool)
		// ------
		// --- declaring totalBytesRead ---
		totalBytesRead := uint64(0)
		// ------
		// --- declaring list of latencies
		latency := make([]int64, 1)
		// ------
		return &HBN{
				url: url,
				method: method,
				duration: duration,
				workers: workers,
				headers: headers,
				useragents: useragents,
				client: &http.Client{},
				s: s,
				errors: 0,
				succeses: 0,
				totalBytesRead: totalBytesRead,
				minLatency: 0.1,
				maxLatency: 0.1,
				latency: latency,
		}, nil
}

// function that user start
func (d *HBN) Run() {
		// starting new bar
		bar := pb.StartNew(d.duration)
		// start testing
		d.start()
		// hello message
		fmt.Printf("Running %ds test for %s\n", d.duration, d.url)
		// bar asynchronous function
		go func() {
				for i := 0; i < d.duration; i++ {
						bar.Increment()
						time.Sleep(time.Second)
				}
		}()
		// waiting for n seconds to end testong
		time.Sleep(time.Duration(d.duration)*time.Second)
		// finish bar working
		bar.Finish()
		// stop testing
		d.stop()
		// --- statistics
		fmt.Printf("total bytes read: %db\n", d.totalBytesRead)
		fmt.Printf("total errors: %d, total succes requests: %d\n", d.errors, d.succeses)
		// ------
		// finding avarage latency
		avrgLatency := findAvrgLatency(d.latency)
		// converting to seconds
		avrg := float32(avrgLatency)/float32(1e9)
		fmt.Printf("Avarage latency is %fs\n", avrg)
		// calculating min and max in list of latencies
		min, max := MinMax(d.latency[1:])
		// converting nanoseconds to seconds
		d.minLatency = convert_nanoseconds_to_seconds(min)
		// converting nanoseconds to seconds
		d.maxLatency = convert_nanoseconds_to_seconds(max)
		// printing minimal and maximum latencies
		fmt.Printf("minimal latancy is: %fs\nmaximum latency is: %fs\n", d.minLatency, d.maxLatency)
		// tnx function
		fmt.Println("tnx for using my tool)))")
}

// start DOSing
func (d *HBN) start() {
		// pull into stop channel flase statement
		go func(){d.s<-false}()
		// starting workers
		for i := 0; i < d.workers; i++ {
				go func() {
						// generating seed to provide the real random number
						rand.Seed(time.Now().UnixNano())
						// random useragent
						useragent := d.useragents[rand.Intn(len(d.useragents))]
						// starting infinity loop testing
						for {
								select {
								// if stop channel is true - end testing
								case <- d.s:
										return
								default:
										// if stop channel is true - testing
										// start main sttack function
										err, p := d.attack(useragent, d.headers)
										// checking for error
										if err != nil {
												// checking if error must to be printed
												if p {
														fmt.Println(err)
														os.Exit(1)
												} else {
														// if error mustn't to be printed - jsut add count of erros
														d.errors += 1
												}
										} else {
												// if wasn't any error - add count of succes requests
												d.succeses += 1
										}
						}
				}
				}()
		}
		// cleaning garbage
		runtime.Gosched()
}

// stop DOSing
func (d *HBN) stop() {
		// pull into stop channel true statement
		go func(){d.s<-true}()
}

// making http request
func (d *HBN) attack(useragent string, headers map[string]string) (error, bool) {
		// if method is GET - making get request
		if d.method == "GET" {
				req, err := http.NewRequest(d.method, d.url, nil)
				if err != nil {
						// error, mustn't to be printed
						return err, false
				} else {
						// set useragent
						req.Header.Set("User-Agent", useragent)
						// setting headers from config file
						for key, value := range headers {
								req.Header.Set(key, value)
						}
						// start time
						start := time.Now().UnixNano()
						// do request
						resp, err := d.client.Do(req)
						// end time
						end := time.Now().UnixNano()
						d.latency = append(d.latency, (end-start))
						//fmt.Println(end.Sub(start).Round(time.Millisecond))
						if err != nil {
								// error mustn't to be printed
								return err, false
						}
						// new bytes read to count total read bytes
						bytesRead, _ := ioutil.ReadAll(resp.Body)
						// adding total read bytes to counter
						atomic.AddUint64(&d.totalBytesRead, uint64(len(bytesRead)))
						// close body
						resp.Body.Close()
						// no error mustn't to be printed
						return nil, false
				}
		} else {
				// if method isn't allowed - return error which must to be printed
				return ErrMethodDoesNotAllowed(d.method), true
		}
}
