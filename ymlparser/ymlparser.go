//Package ymlparser implements structs for yml parsing
package ymlparser

import (
	"io/ioutil"
	"log"
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
	Services []Service
	VMs      []VM
	Name     string
	ISODate  time.Time
	timer    *time.Timer
}

//Config provides data of the cloud infrastructure.
type Config struct {
	States   []State `yaml:"states"`
	Provider struct {
		Name        string
		ManagerIP   string `yaml:"managerIP"`
		Username    string
		Password    string
		ConfigFile  string `yaml:"configFile"`
		JoinCommand string `yaml:"joinCommand"`
	} `yaml:"provider"`
}

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
	log.Printf("%s parsed correctly", configFile)

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

func (s State) getReadableTime() string {
	return time.Now().Format(time.RFC822)
}
