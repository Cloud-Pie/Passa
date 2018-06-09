package main

//go:generate go run internal/generate/ymlGenerator.go

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Cloud-Pie/Passa/cloudsolution/dockerswarm"
	"github.com/Cloud-Pie/Passa/cloudsolution/lrz"
	"github.com/Cloud-Pie/Passa/database"

	"github.com/Cloud-Pie/Passa/cloudsolution"
	"github.com/Cloud-Pie/Passa/notification"
	"github.com/Cloud-Pie/Passa/notification/consoleprinter"
	"github.com/Cloud-Pie/Passa/notification/telegram"
	"github.com/Cloud-Pie/Passa/server"
	"github.com/Cloud-Pie/Passa/ymlparser"
)

const (
	defaultLogFile       = "test.log"
	defaultYMLFile       = "test/passa-states-test.yml"
	defaultCheckInterval = 20 //secs
)

var notifier notification.NotifierInterface
var flagVars flagVariable

func main() {
	database.InitializeDB()
	var err error
	stateChannel := make(chan *ymlparser.State) //Communication between server states and our scheduler
	flagVars = parseFlags()
	c := ymlparser.ParseStatesfile(flagVars.configFile)

	//Notifier code Start
	notifier, err = telegram.InitializeClient()

	if err != nil {
		notifier = consoleprinter.InitializeClient()
	}
	//Notifier code End

	//Code For Cloud Management Start

	var cloudManager cloudsolution.CloudManagerInterface
	if !flagVars.noCloud {
		if c.Provider.Name == "docker-swarm" {
			cloudManager = dockerswarm.NewSwarmManager(c.Provider.ManagerIP)
		} else if c.Provider.Name == "lrz" {
			cloudManager = lrz.NewLRZManager(c.Provider.Username, c.Provider.Password, c.Provider.ConfigFile)
		}
	}

	go schedulerRoutine(stateChannel, cloudManager) //Most important line in whole project

	for idx := range c.States {

		if !c.States[idx].ISODate.IsZero() {
			stateChannel <- &c.States[idx]
		} else {
			log.Printf("Invalid time for: %s", c.States[idx].Name)
		}
	}
	//Code For Cloud Management End

	//Code for WatchDog Start
	if !flagVars.noCloud {
		go periodicCheckRoutine(cloudManager)
	}
	//Code for WatchDog End
	//Server code Start
	server := server.SetupServer(c, stateChannel)
	server.Run()
	//Server code End
}

//Most important function in the whole project !!
func schedulerRoutine(stateChannel chan *ymlparser.State, cm cloudsolution.CloudManagerInterface) {
	for incomingState := range stateChannel {
		durationUntilStateChange := incomingState.ISODate.Sub(time.Now())
		deploymentTimer := time.AfterFunc(durationUntilStateChange, scale(cm, *incomingState)) //Golang closures
		incomingState.SetTimer(deploymentTimer)
		fmt.Printf("Saved Deployment: %v\n", incomingState)
		database.InsertState(*incomingState)

	}
}

func scale(manager cloudsolution.CloudManagerInterface, s ymlparser.State) func() {

	return func() {

		if !flagVars.noCloud {
			manager = manager.ChangeState(s)
			fmt.Printf("%#v", manager.GetLastDeployedState())
		}
		notifier.Notify("Deployed " + s.Name)

	}
}

func setLogFile(lf string) string {
	if lf == "" {
		lf = defaultLogFile
	}
	fmt.Println("Writing log to  -> ", lf)
	os.MkdirAll(filepath.Dir(lf), 0700)
	f, err := os.OpenFile(lf, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	return lf
}

type flagVariable struct {
	noCloud             bool
	configFile          string
	logFile             string
	checkStatusInterval int
}

func parseFlags() flagVariable {
	noCloud := flag.Bool("no-cloud", false, "Don't start cloud management") //NOTE: For testing only
	configFile := flag.String("state-file", defaultYMLFile, "config file")
	logFile := flag.String("test-file", defaultLogFile, "log file")
	statusInterval := flag.Int("check-status-interval", defaultCheckInterval, "Check every <this> second")

	flag.Parse()
	return flagVariable{
		noCloud:             *noCloud,
		configFile:          *configFile,
		logFile:             *logFile,
		checkStatusInterval: *statusInterval,
	}
}

func periodicCheckRoutine(cm cloudsolution.CloudManagerInterface) {

	sleepDuration := time.Duration(flagVars.checkStatusInterval) * time.Second
	for ; ; time.Sleep(sleepDuration) {
		if !cm.CheckState() {
			log.Println("False State")
		}
	}
}
