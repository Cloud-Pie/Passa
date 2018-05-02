package ymlparser

import (
	"fmt"
	"testing"
	"time"
)

func Test_parseTime(t *testing.T) {
	c := ParseStatesfile("../test/passa-states-test.yml")

	tis, err := time.Parse(TimeLayout, c.MyTime)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tis)
	fmt.Println(time.Now())
	duration := tis.Sub(time.Now())
	fmt.Println(duration.Seconds())
}
