//Package ymlparser implements structs for yml parsing
package ymlparser

import (
	"io/ioutil"

	"time"

	logging "github.com/op/go-logging"
	yaml "gopkg.in/yaml.v2"
)

var log = logging.MustGetLogger("passa")

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
	ID       string
	Services []Service
	VMs      []VM
	Name     string
	ISODate  time.Time
	Timer    *time.Timer
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
	log.Debug("%s parse correctly", configFile)

	return c
}

func (s State) getReadableTime() string {
	return time.Now().Format(time.RFC822)
}
