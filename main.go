package main

//go:generate go run internal/generate/ymlGenerator.go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gitlab.lrz.de/ga53lis/PASSA/notifier"
	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
)

const (
	defaultLogFile = "test.log"
	defaultYMLFile = "passa-states.yml"
)

func main() {
	var wg sync.WaitGroup
	c := ymlparser.ParseStatesfile(defaultYMLFile)
	notifier.InitializeClient() //FIXME: this will definitely change
	notifier.Notify("Connected to PASSA")

	wg.Add(len(c.States))
	currentTime := time.Now()

	for _, state := range c.States {

		predictedTime, err := time.Parse(ymlparser.TimeLayout, state.Time)
		if err != nil {
			panic(err)
		}
		durationUntilStateChange := predictedTime.Sub(currentTime)
		time.AfterFunc(durationUntilStateChange, scale(state, &wg)) //Golang closures
	}

	wg.Wait()
	fmt.Println("Exiting")
}

func scale(s ymlparser.State, wg *sync.WaitGroup) func() {

	return func() {
		defer wg.Done()
		log.Println(s.Name)
		//ensureState(s)
		notifier.Notify("Deploying " + s.Name)
	}
}

func getCurrentServices(p string) ([]ymlparser.Service, error) {
	var currentServices []ymlparser.Service
	resp, err := http.Get(p + "/status")
	if err != nil {
		// handle err
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	json.Unmarshal(body, &currentServices)
	defer resp.Body.Close()
	return currentServices, nil
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
