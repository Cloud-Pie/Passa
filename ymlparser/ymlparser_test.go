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
