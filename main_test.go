package main

//run go generate first
import (
	"fmt"
	"testing"
	"time"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

func TestChannel(t *testing.T) {
	myChan := make(chan *ymlparser.State, 1)

	myState := ymlparser.State{
		Name:     "zzzz",
		Services: []ymlparser.Service{{Name: "relax_web", Scale: 1}, {Name: "relax_visualizer", Scale: 1}},
		VMs:      []ymlparser.VM{{Type: "asd", Scale: 1}},
	}

	myChan <- &myState
	newState := <-myChan
	newState.Name = "aaa"

	fmt.Printf("%v", myState)
}

func Test_minusDuration(t *testing.T) {
	myDuration := time.Now().Local().Add(time.Hour * -2)

	if time.Now().After(myDuration) {
		fmt.Println("past time")
	}
}
