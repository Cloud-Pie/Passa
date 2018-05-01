package main

//run go generate first
import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
)

func Test_parseTime(t *testing.T) {
	c := ymlparser.ParseStatesfile("passa-states.yml")
	layout := "02-01-2006, 15:04:05 MST" //GOLANG's special time thing, cost me 20 mins

	tis, err := time.Parse(layout, c.MyTime)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tis)
	fmt.Println(time.Now())
	duration := tis.Sub(time.Now())
	fmt.Println(duration.Seconds())
}
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

/*func Test_getCurrentService(t *testing.T) {
	const providerURL = "http://localhost:4000"
	currentServices, err := getCurrentServices(providerURL)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%v", currentServices)
}*/
