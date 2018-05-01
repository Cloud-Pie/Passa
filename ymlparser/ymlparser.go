package ymlparser

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Service struct {
	Name string `json:"Name"`
	Port int    `json:"Port"`
}
type State struct {
	Time     string
	Services []Service
	Name     string
}
type config struct {
	Version     string  `yaml:"version"`
	States      []State `yaml:"states"`
	ProviderURL string  `yaml:"providerURL"`
	MyTime      string  `yaml:"myTime"`
}

var providerURL string

//TimeLayout is the golang's special time format
const TimeLayout = "02-01-2006, 15:04:05 MST"

func ParseStatesfile(configFile string) *config {
	var c *config
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
