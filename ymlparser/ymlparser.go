//Package ymlparser implements structs for yml parsing
package ymlparser

import (
	"io/ioutil"

	"time"

	logging "github.com/op/go-logging"
	yaml "gopkg.in/yaml.v2"
)

var log = logging.MustGetLogger("passa")

//Service keeps cpu and memory for one service
type Service map[string]ServiceInfo

//ServiceInfo is constraints for kubernetes
type ServiceInfo struct {
	CPU      string
	Memory   string
	Replicas int
}

//VM keeps the type and scale of virtual machines.
type VM map[string]int

//State is the metadata of the state expected to scale to.
type State struct {
	Services     Service
	VMs          VM
	Name         string
	ISODate      time.Time
	timer        *time.Timer
	ExpectedTime time.Time
	RealTime     time.Time
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
		ClusterName string `yaml:"clusterName"`
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
	log.Debug("%s parsed correctly", configFile)

	return c
}

func (s State) getReadableTime() string {
	return time.Now().Format(time.RFC822)
}

//SetTimer sets the timer of the state
func (s *State) SetTimer(t *time.Timer) {
	s.timer = t
}

//StopTimer stops the timer of the state
func (s *State) StopTimer() {
	if s.timer != nil {
		s.timer.Stop()
	}
}
