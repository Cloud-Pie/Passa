package lrz

import (
	"fmt"
	"testing"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

func Test_lrz_getActiveState(t *testing.T) {
	tester := NewLRZManager("username", "password", "/Users/atakanyenel/Desktop/Computer_Science/go/src/github.com/Cloud-Pie/Passa/test/admin.conf", "")

	state := tester.GetActiveState()
	fmt.Printf("%+v", state)

}

func Test_lrz_scaleContainers(t *testing.T) {
	tester := NewLRZManager("username", "password", "/Users/atakanyenel/Desktop/Computer_Science/go/src/github.com/Cloud-Pie/Passa/test/admin.conf", "")

	state := tester.GetActiveState()

	fmt.Printf("%+v", state)

	wantedState := ymlparser.State{Services: ymlparser.Service{
		"hello_world": ymlparser.ServiceInfo{Replicas: 4},
		"nginx":       ymlparser.ServiceInfo{Replicas: 3},
	}}

	tester.ChangeState(wantedState)

}
