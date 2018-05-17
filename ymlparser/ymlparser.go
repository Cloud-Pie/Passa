//Package ymlparser implements structs for yml parsing
package ymlparser

import (
	"fmt"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

//Service keeps the name and scale of the scaled service.
type Service struct {
	Name  string
	Scale int
}

//VM keeps the type and scale of virtual machines.
type VM struct {
	Type  string
	Scale int
}

//State is the metadata of the state expected to scale to.
type State struct {
	Time     string
	Services []Service
	VMs      []VM
	Name     string
	ISODate  time.Time
	timer    *time.Timer
}

//Config provides data of the cloud infrastructure.
type Config struct {
	Version  string  `yaml:"version"`
	States   []State `yaml:"states"`
	Provider struct {
		Name      string
		ManagerIP string `yaml:"managerIP"`
	} `yaml:"provider"`
}

var providerURL string

//TimeLayout is the golang's special time format
const TimeLayout = "02-01-2006, 15:04:05 MST"

//ParseStatesfile parses the states file according to configuration.
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

	for idx := range c.States {
		isoTimeFormat, err := time.Parse(TimeLayout, c.States[idx].Time)
		if err != nil {
			panic(err)
		}
		c.States[idx].ISODate = isoTimeFormat
	}

	return c
}

//SetTimer sets the private timer variable
func (s *State) SetTimer(t *time.Timer) {
	if s.timer != nil {
		dur := s.ISODate.Sub(time.Now())
		s.timer.Reset(dur)
	}
	s.timer = t
}
