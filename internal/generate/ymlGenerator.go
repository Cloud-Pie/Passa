package main

//check test directory
import (
	"fmt"
	"os"
	"time"

	"github.com/Cloud-Pie/Passa/ymlparser"
	"gopkg.in/yaml.v2"
)

const ymlReference = "passa-states.yml"
const ymlTest = "test/passa-states-test.yml"

func main() {
	c := ymlparser.ParseStatesfile(ymlReference)
	currentTime := time.Now()

	addedMinutes := [5]int{-2, 10, 15, 20, 25} //Constant

	for idx := range c.States {
		timein := currentTime.Local().Add(time.Hour * 24 * 30 * time.Duration(addedMinutes[idx]))
		c.States[idx].ISODate = timein
	}
	fmt.Printf("%v", c)

	ymlByte, err := yaml.Marshal(&c)
	check(err)
	f, err := os.Create(ymlTest)
	check(err)
	defer f.Close()

	f.WriteString("# Generated by internal/generator/ymlGenerator.go\n---\n")
	f.Write(ymlByte)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
