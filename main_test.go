package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	yaml "gopkg.in/yaml.v2"
)

func Test_parseTime(t *testing.T) {
	var c config
	statesFile := flag.String("states", defaultYMLFile, "path of yml file")
	fmt.Println(*statesFile)
	source, err := ioutil.ReadFile(*statesFile)
	if err != nil {
		panic(err)
	}
	_ = yaml.Unmarshal(source, &c)

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

func Test_readYML(t *testing.T) {
	var c config
	source, err := ioutil.ReadFile(defaultYMLFile)
	if err != nil {
		panic(err)
	}
	_ = yaml.Unmarshal(source, &c)

	fmt.Printf("Version: %v\n", c.Version)
	providerURL = c.ProviderURL
	fmt.Printf("ProviderURL: %s\n", providerURL)
}

/*func Test_getCurrentService(t *testing.T) {
	const providerURL = "http://localhost:4000"
	currentServices, err := getCurrentServices(providerURL)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%v", currentServices)
}*/
