package main

//go:generate go run internal/generate/ymlGenerator.go

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Cloud-Pie/Passa/cloudsolution/dockerswarm"
	"github.com/Cloud-Pie/Passa/cloudsolution/lrz"
	"github.com/Cloud-Pie/Passa/database"

	"github.com/Cloud-Pie/Passa/cloudsolution"
	"github.com/Cloud-Pie/Passa/notification"
	"github.com/Cloud-Pie/Passa/notification/consoleprinter"
	"github.com/Cloud-Pie/Passa/notification/telegram"
	"github.com/Cloud-Pie/Passa/server"
	"github.com/Cloud-Pie/Passa/ymlparser"
	"github.com/op/go-logging"
)

const (
	defaultLogFile       = "test.log"
	defaultYMLFile       = "passa-states.yml"
	defaultCheckInterval = 20 //secs
)

var notifier notification.NotifierInterface
var flagVars flagVariable
var cloudManager cloudsolution.CloudManagerInterface

var isCurrentlyDeploying bool
var isStateTrue bool
var log = logging.MustGetLogger("passa")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	styleEntry()
	var err error
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

	stateChannel := make(chan *ymlparser.State) //Communication between server states and our scheduler
	flagVars = parseFlags()
	c := ymlparser.ParseStatesfile(flagVars.configFile)
	database.InitializeDB()
	//Notifier code Start
	notifier, err = telegram.InitializeClient()

	if err != nil {
		log.Debug("Can't connect telegram")
		log.Info("console")
		notifier = consoleprinter.InitializeClient()
	} else {
		log.Info("telegram")
		compName, _ := os.Hostname()
		notifier.Notify(fmt.Sprintf("Passa connected to telegram at %s", compName))
	}
	//Notifier code End

	//Code For Cloud Management Start

	if !flagVars.noCloud {
		if c.Provider.Name == "docker-swarm" {
			cloudManager = dockerswarm.NewSwarmManager(c.Provider.ManagerIP)
		} else if c.Provider.Name == "lrz" {
			cloudManager = lrz.NewLRZManager(c.Provider.Username, c.Provider.Password, c.Provider.ConfigFile, c.Provider.JoinCommand)
		}
	}

	go schedulerRoutine(stateChannel, cloudManager) //Most important line in whole project

	for idx := range c.States {

		if !c.States[idx].ISODate.IsZero() {
			stateChannel <- &c.States[idx]
		} else {
			log.Warning("Invalid time for: %s", c.States[idx].Name)
		}
	}
	//Code For Cloud Management End

	//Code for WatchDog Start
	if !flagVars.noCloud {
		go periodicCheckRoutine()
	}
	//Code for WatchDog End
	//Server code Start
	server := server.SetupServer(stateChannel)
	server.GET("/api/current", func(c *gin.Context) {
		c.JSON(200,
			gin.H{
				"lastDeployed": cloudManager.GetLastDeployedState(),
				"active":       cloudManager.GetActiveState(),
				"isStateTrue":  isStateTrue,
			})
	})
	log.Info("Server listening on port 5555")
	server.Run(":5555")
	//Server code End
}

//Most important function in the whole project !!
func schedulerRoutine(stateChannel chan *ymlparser.State, cm cloudsolution.CloudManagerInterface) {
	for incomingState := range stateChannel {
		if time.Now().After(incomingState.ISODate) && false { //FIXME: remove && false
			log.Notice("%s is a past state, not deploying\n", incomingState.Name)
			database.InsertState(*incomingState)
		} else {
			durationUntilStateChange := incomingState.ISODate.Sub(time.Now())

			deploymentTimer := time.AfterFunc(durationUntilStateChange, scale(*incomingState)) //Golang closures
			incomingState.Timer = deploymentTimer
			database.InsertState(*incomingState)
		}
	}
}

func scale(s ymlparser.State) func() {

	return func() {

		if !flagVars.noCloud {
			isCurrentlyDeploying = true
			cloudManager = cloudManager.ChangeState(s)
			log.Info("%#v", cloudManager.GetLastDeployedState())
			isCurrentlyDeploying = false

		}
		notifier.Notify("Deployed " + s.Name)

	}
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

func periodicCheckRoutine() {

	sleepDuration := time.Duration(flagVars.checkStatusInterval) * time.Second
	//time.Sleep(sleepDuration)
	for ; ; time.Sleep(sleepDuration) {

		if isCurrentlyDeploying {
			log.Notice("Actively Deploying new State...")
		} else {
			switch isStateTrue = cloudManager.CheckState(); isStateTrue {
			case true:
				log.Info("State checked, everything is fine")
			case false:
				//NOTE: send SPDT notification
				log.Warning("False State")

			}

		}
	}
}

func styleEntry() {
	fmt.Println(`
	.______      ___           _______.     _______.     ___     
	|   _  \    /   \         /       |    /       |    /   \    
	|  |_)  |  /  ^  \       |   (----'   |   (----'   /  ^  \   
	|   ___/  /  /_\  \       \   \        \   \      /  /_\  \  
	|  |     /  _____  \  .----)   |   .----)   |    /  _____  \ 
	| _|    /__/     \__\ |_______/    |_______/    /__/     \__\
	`)
}
