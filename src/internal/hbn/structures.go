package hbn

// Config structure to parse from json file
type Config struct {
	Useragents string   `json:"useragents"`
	Headers    []string `json:"headers"`
}
