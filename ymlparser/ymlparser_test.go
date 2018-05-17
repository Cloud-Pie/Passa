package ymlparser

import (
	"fmt"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	timeString := "10-05-2018, 23:51:50 CEST"
	jsTimeFormat, err := time.Parse(TimeLayout, timeString)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n%s", jsTimeFormat, jsTimeFormat.Format(time.RFC3339))
}

func TestSetTimer(t *testing.T) {
	s := State{
		Name: "testState",
		Time: "18-08-2018, 20:00:00 CEST",
	}

	s.ISODate, _ = time.Parse(TimeLayout, s.Time)
	myTimer := time.AfterFunc(s.ISODate.Sub(time.Now()), func() {
		fmt.Println("something something")
	})
	s.SetTimer(myTimer)
	fmt.Printf("%+v", s)
	s.Time = "18-08-2019, 21:00:00 CEST"
	s.ISODate, _ = time.Parse(TimeLayout, s.Time)
	myNewTimer := time.AfterFunc(s.ISODate.Sub(time.Now()), func() {
		fmt.Println("something something")
	})

	s.SetTimer(myNewTimer)
	fmt.Printf("%+v", s)
}

func TestParseStateFile(t *testing.T) {
	c := ParseStatesfile("../test/passa-states-test.yml")
	if c.Version != "0.6" {
		t.Fail()
	}
	if c.Provider.Name != "docker-swarm" {
		t.Fail()
	}
}
