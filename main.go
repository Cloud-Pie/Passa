package main

//go:generate go run internal/generate/ymlGenerator.go

import (
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

var providerURL string

func main() {
	var wg sync.WaitGroup
	c := ymlparser.ParseStatesfile(defaultYMLFile)
	notifier.InitializeClient() //FIXME: this will definitely change
	notifier.Notify("Connected to PASSA")
	providerURL = c.ProviderURL
	cloudManager := cloudsolution.NewSwarmManager(providerURL)
	wg.Add(len(c.States))
	currentTime := time.Now()

	for _, state := range c.States {

		durationUntilStateChange := state.ISODate.Sub(currentTime)
		time.AfterFunc(durationUntilStateChange, scale(cloudManager, state, &wg)) //Golang closures
	}

	fmt.Println("Exiting")

	//Server start
	server := server.SetupServer(c)
	server.Run()
	//So the program doesn't end
	wg.Wait() //TODO: maybe we can remove this all together.

}

func scale(manager cloudsolution.CloudManager, s ymlparser.State, wg *sync.WaitGroup) func() {

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
