package main

//go:generate go run internal/generate/ymlGenerator.go

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gitlab.lrz.de/ga53lis/PASSA/cloudsolution"
	"gitlab.lrz.de/ga53lis/PASSA/notifier"
	"gitlab.lrz.de/ga53lis/PASSA/server"
	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
)

const (
	defaultLogFile = "test.log"
	defaultYMLFile = "test/passa-states-test.yml"
)

func main() {

	flagVars := parseFlags()
	var wg sync.WaitGroup
	c := ymlparser.ParseStatesfile(flagVars.configFile)

	//Notifier code Start
	notifier.InitializeClient() //FIXME: this will definitely change
	notifier.Notify("Connected to PASSA")
	//Notifier code End

	//Code For Cloud Management Start

	if !flagVars.noCloud {
		cloudManager := cloudsolution.NewSwarmManager(c.ProviderURL)

		for _, state := range c.States {

			durationUntilStateChange := state.ISODate.Sub(time.Now())

			if durationUntilStateChange < 0 {
				log.Printf("Duration already passed:%v", durationUntilStateChange)
			} else {
				time.AfterFunc(durationUntilStateChange, scale(cloudManager, state, &wg)) //Golang closures
				wg.Add(1)
			}
		}
		//Code For Cloud Management End
	}
	//Server code Start
	server := server.SetupServer(c)
	server.Run()
	//Server code End

	wg.Wait() //TODO: maybe we can remove this all together.
}

func scale(manager cloudsolution.CloudManager /*Interface*/, s ymlparser.State, wg *sync.WaitGroup) func() {

	return func() {
		defer wg.Done()
		for _, service := range s.Services {

			fmt.Println(manager.ChangeState(service))

			notifier.Notify("Deployed " + s.Name)
		}
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
	noCloud    bool
	configFile string
	logFile    string
}

func parseFlags() flagVariable {
	noCloud := flag.Bool("no-cloud", false, "Don't start cloud management")
	configFile := flag.String("state-file", defaultYMLFile, "config file")
	logFile := flag.String("test-file", defaultLogFile, "log file")

	flag.Parse()
	return flagVariable{
		noCloud:    *noCloud,
		configFile: *configFile,
		logFile:    *logFile,
	}
}
