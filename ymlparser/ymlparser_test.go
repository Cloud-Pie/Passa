package ymlparser

import (
	"fmt"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {

	jsTimeFormat := time.Now().Format(time.RFC822)

	fmt.Printf("%s\n", jsTimeFormat)
}

func TestSetTimer(t *testing.T) {
	s := State{
		Name:    "testState",
		ISODate: time.Now().Add(time.Hour * time.Duration(2)), //after 2 hours
	}

	myTimer := time.AfterFunc(s.ISODate.Sub(time.Now()), func() {
		fmt.Println("something something")
	})
	s.SetTimer(myTimer)
	fmt.Printf("%+v", s)
	s.ISODate = time.Now().Add(time.Hour * 24 * 365 * time.Duration(1))
	myNewTimer := time.AfterFunc(s.ISODate.Sub(time.Now()), func() {
		fmt.Println("something something")
	})

	s.SetTimer(myNewTimer)
	fmt.Printf("%+v", s)
}

func TestParseStateFile(t *testing.T) {
	ParseStatesfile("../test/passa-states-test.yml")
}

func TestState_getReadableTime(t *testing.T) {
	type fields struct {
		Services []Service
		VMs      []VM
		Name     string
		ISODate  time.Time
		timer    *time.Timer
	}
	myTime := time.Now()
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Valid Time",
			fields: fields{
				ISODate: myTime,
			},
			want: myTime.Format(time.RFC822),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := State{

				ISODate: tt.fields.ISODate,
			}
			if got := s.getReadableTime(); got != tt.want {
				t.Errorf("State.getReadableTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
