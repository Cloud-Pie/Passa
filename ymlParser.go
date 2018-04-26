package main

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type service struct {
	Name string `json:"Name"`
	Port int    `json:"Port"`
}
type state struct {
	Time     string
	Services []service
	Name     string
}
type config struct {
	Version     string  `yaml:"version"`
	States      []state `yaml:"states"`
	ProviderURL string  `yaml:"providerURL"`
	MyTime      string  `yaml:"myTime"`
}

var providerURL string

const timeLayout = "02-01-2006, 15:04:05 MST"

func parseStatesfile(configFile string) *config {
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
