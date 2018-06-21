package database

import (
	"encoding/json"
	"fmt"

	"github.com/op/go-logging"

	"github.com/Cloud-Pie/Passa/ymlparser"

	"github.com/nanobox-io/golang-scribble"
)

var db *scribble.Driver
var log = logging.MustGetLogger("passa")

const dbName = "state"
const dir = "./.db/"

//InitializeDB initializes the database
func InitializeDB() {
	var err error
	db, err = scribble.New(dir, nil)
	if err != nil {
		log.Error("Database not initialized...")
	}
}

//InsertState inserts state to DB
func InsertState(newState ymlparser.State) {
	db.Write(dbName, newState.ID, newState)
}

//GetSingleState returns single state
func GetSingleState(stateID string) ymlparser.State {

	state := ymlparser.State{}
	if err := db.Read(dbName, stateID, &state); err != nil {
		log.Warning("Couldn't get %s", stateID)
		return ymlparser.State{}
	}
	return state

}

//ReadAllStates returns all states
func ReadAllStates() []ymlparser.State {
	records, err := db.ReadAll(dbName)
	if err != nil {
		log.Error("Error", err)
		return nil
	}
	returnStates := []ymlparser.State{}
	for _, f := range records {
		stateFound := ymlparser.State{}
		if err := json.Unmarshal([]byte(f), &stateFound); err != nil {
			fmt.Println("Error", err)
		}
		returnStates = append(returnStates, stateFound)
	}
	return returnStates
}

//DeleteState deletes state with that name
func DeleteState(stateID string) error {
	var err error

	GetSingleState(stateID).Timer.Stop()

	if err = db.Delete(dbName, stateID); err != nil {
		fmt.Println("Error", err)

	}
	return err
}

//UpdateState updates state
func UpdateState(newState ymlparser.State, oldStateID string) {

	DeleteState(oldStateID)
	InsertState(newState)

}
