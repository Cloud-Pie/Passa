//Package database provides functions for database
package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

type myDataBase struct {
	filepath string
	sync.Mutex
}

var db myDataBase

const fileName = ".db.json"

//InitializeDB initializes a new no-sql Database.
func InitializeDB() {

	//goPath := os.Getenv("GOPATH")
	//filePath := goPath + "/src/github.com/Cloud-Pie/Passa/database/" + fileName
	db = myDataBase{
		filepath: fileName,
	}
	var file *os.File
	if _, err := os.Stat(db.filepath); os.IsNotExist(err) {
		file, err = os.Create(db.filepath)
		if err != nil {
			panic(err)
		}
		dbByte, err := json.Marshal(&ymlparser.Config{})
		if err != nil {
			panic(err)
		}
		file.Write(dbByte)
		defer file.Close()
	}

}

//InsertState inserts state in to the database
func InsertState(newState ymlparser.State) {

	db.Lock()
	defer db.Unlock()

	c := loadDBtoMemory()
	c.States = append(c.States, newState)

	writeToFile(c)
}

//ReadAllStates reads all the states from the database
func ReadAllStates() []ymlparser.State {

	c := loadDBtoMemory()
	return c.States

}

//SearchQuery returns the index of the state in config file
func searchQuery(currentStates []ymlparser.State, searchName string) int {

	for idx := range currentStates {
		if currentStates[idx].Name == searchName {
			return idx
		}
	}
	return -1
}

func loadDBtoMemory() ymlparser.Config {

	if db.filepath == "" {
		panic(errors.New("No DB initialized"))
	}
	var c ymlparser.Config

	source, err := ioutil.ReadFile(db.filepath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(source, &c)
	if err != nil {
		panic(err)
	}

	return c

}

//UpdateState updates the state with the given name
func UpdateState(updatedstate ymlparser.State, oldStateName string) error {
	db.Lock()
	defer db.Unlock()
	c := loadDBtoMemory()
	pos := searchQuery(c.States, oldStateName)
	if pos == -1 {
		return errors.New("No state with that name")
	}
	c.States[pos] = updatedstate
	writeToFile(c)
	return nil
}

//DeleteState deletes the state with the given name
func DeleteState(stateToDelete string) error {
	db.Lock()
	defer db.Unlock()
	c := loadDBtoMemory()
	pos := searchQuery(c.States, stateToDelete)
	if pos == -1 {
		return errors.New("No state with that name")
	}
	c.States[pos] = c.States[len(c.States)-1]
	c.States[len(c.States)-1] = ymlparser.State{}
	c.States = c.States[:len(c.States)-1]
	writeToFile(c)
	return nil
}
func dropDB() { //Just for testing
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return
	}
	fmt.Println(db.filepath)
	var err = os.Remove(fileName)
	if err != nil {
		panic(err)
	}

}

func writeToFile(c ymlparser.Config) {
	dbByte, err := json.Marshal(&c)
	if err != nil {
		panic(err)
	}

	f, _ := os.OpenFile(db.filepath, os.O_RDWR, 0644)
	defer f.Close()
	f.Write(dbByte)
}

//GetSingleState gets the state with the given name
func GetSingleState(stateName string) (ymlparser.State, error) {
	c := loadDBtoMemory()
	pos := searchQuery(c.States, stateName)
	if pos == -1 {
		return ymlparser.State{}, errors.New("No state found")
	}
	return c.States[pos], nil
}
