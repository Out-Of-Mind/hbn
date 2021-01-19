package hbn

import (
		"os"
		"fmt"
		"bufio"
		"strings"
		"encoding/json"
		errs "github.com/out-of-mind/hbn/src/internal/errors"
)

func findAvrgLatency(list []int64) int64 {
		res := int64(0)
		for _, i := range list[1:] {
				res += i
		}
		res /= int64(len(list[1:]))
		return res
}

func MinMax(values []int64) (int64, int64) {
		min_value := values[0]
		max_value := values[0]
		for _, i := range values {
				if i < min_value {
						min_value = i
				} else if i > max_value {
						max_value = i
				}
		}

		return min_value, max_value
}

func convert_nanoseconds_to_seconds(value int64) float32 {
		res := float32(value)/float32(1e9)
		return res
}

func getConfig(path_to_config string) *Config {
		f, err := os.Open(path_to_config)
		if err != nil {
				fmt.Println(errs.ErrConfigWasNotFound(path_to_config))
				os.Exit(1)
		}
		decoder := json.NewDecoder(f)
		config := new(Config)
		err = decoder.Decode(&config)
		if err != nil {
				fmt.Println(errs.ErrConfigRead())
				os.Exit(1)
		}
		f.Close()

		return config
}

func getUseragents(use_useragents bool, path_to_useragents string) []string {
		useragents := make([]string, 1)
		if use_useragents {
				f, err := os.Open(path_to_useragents)
				if err != nil {
						fmt.Println(errs.ErrUseragentsWasNotFound(path_to_useragents))
						os.Exit(1)
				}
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
						useragents = append(useragents, scanner.Text())
				}
				if err := scanner.Err(); err != nil {
						fmt.Println(errs.ErrUseragentsFileRead())
						os.Exit(1)
				}
				f.Close()

				return useragents
		} else {
				return useragents
		}
}

func getHeaders(use_headers bool, headers_unsplitted []string) map[string]string {
		headers := make(map[string]string)
		if use_headers {
				for _, v := range headers_unsplitted {
						o := strings.Split(v, " ")
						headers[o[0]] = o[1]
				}
				return headers
		} else {
				return headers
		}
}

func getCookies(use_cookies bool, cookies []string) []string {
		if use_cookies {
				return cookies
		} else {
				return make([]string, 1)
		}
}
