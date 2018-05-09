//Package ymlparser implements structs for yml parsing
package ymlparser

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Service struct {
	Name  string `json:"Name"`
	Scale string `json:"Scale"`
}
type State struct {
	Time     string
	Services []Service
	Name     string
}
type Config struct {
	Version     string  `yaml:"version"`
	States      []State `yaml:"states"`
	ProviderURL string  `yaml:"providerURL"`
	MyTime      string  `yaml:"myTime"`
}

var providerURL string

//TimeLayout is the golang's special time format
const TimeLayout = "02-01-2006, 15:04:05 MST"

func ParseStatesfile(configFile string) *Config {
	var c *Config
	source, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Version %v\n", c.Version)

	return c
}
