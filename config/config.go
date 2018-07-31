// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	Url string `config:"url"`
	Authorization string `config:"authorization"`
	//Location string `confing:"location"`
	Method string `config:"method"`
	JsonDotMode string `config:"jsonDotMode"`
	OutputFormat string `config:"outputFormat"`
	DefaultOutputFormat string `config:"DefaultOutputFormat"`
	Headers map[string]string `config:"headers"`
	Fields map[string]string `config:"fields"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}
