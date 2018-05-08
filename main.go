package main

//go:generate go run internal/generate/ymlGenerator.go

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.lrz.de/ga53lis/PASSA/cloudsolution"
	"gitlab.lrz.de/ga53lis/PASSA/notifier"
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
	r := gin.Default()
	notifier.InitializeClient() //FIXME: this will definitely change
	notifier.Notify("Connected to PASSA")
	providerURL = c.ProviderURL
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

	go wg.Wait() //TODO: maybe we can remove this all together.
	fmt.Println("Exiting")

	//Server start
	r.GET("/ui/states", func(cxt *gin.Context) {
		cxt.JSON(200, c)
	})

	r.Run()
}

func scale(s ymlparser.State, wg *sync.WaitGroup) func() {

	return func() {
		defer wg.Done()
		for _, service := range s.Services {
			scaleInt, err := strconv.Atoi(service.Scale)
			if err != nil {
				panic(err)
			}
			fmt.Println(cloudsolution.ChangeState(scaleInt))
			cloudsolution.ScaleContainers(providerURL, service.Name, service.Scale)

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
