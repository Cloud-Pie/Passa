package lrz

import (
	"fmt"
	"testing"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

func Test_lrz_getServiceCount(t *testing.T) {
	tester := NewLRZManager("username", "password", "/Users/atakanyenel/Desktop/Computer_Science/go/src/github.com/Cloud-Pie/Passa/test/kubeconfig.txt")

	services := tester.getServiceCount()

	state := tester.GetActiveState()
	fmt.Printf("%+v", services)
	fmt.Printf("%+v", state)

}

func Test_lrz_scaleContainers(t *testing.T) {
	tester := NewLRZManager("username", "password", "/Users/atakanyenel/Desktop/mycube.txt")

	state := tester.GetActiveState()

	fmt.Printf("%+v", state)

	wantedState := ymlparser.State{Services: []ymlparser.Service{
		ymlparser.Service{
			Name:  "hello-world",
			Scale: 2,
		},
		ymlparser.Service{
			Name:  "my-nginx",
			Scale: 4,
		},
	}}

	tester.ChangeState(wantedState)

}
