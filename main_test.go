package main

//run go generate first
import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

func Test_setLogFile(t *testing.T) {
	type args struct {
		lf string
	}
	tests := []struct {
		name string
		args args
	}{
		{"log in another folder", args{"log/_test.log"}},
		{"log in this folder", args{"_test.log"}},
		{"log with empty string", args{""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileName := setLogFile(tt.args.lf)
			if _, err := os.Stat(fileName); os.IsNotExist(err) {
				t.Fail()
			} else {
				os.Remove(fileName)
				os.Remove(filepath.Dir(fileName))
			}
		})
	}
}

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
