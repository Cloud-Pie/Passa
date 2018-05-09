package server

import (
	"fmt"
	"testing"

	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
)

func TestStartServer(t *testing.T) {
	c := ymlparser.ParseStatesfile("../test/passa-states-test.yml")
	r := StartServer(c)
	for _, rt := range r.Routes() {
		fmt.Printf("%v", rt.Path)
	}
	fmt.Printf("%s", r.Routes())
}
