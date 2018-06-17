package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

/*
	statesRest.POST("/", createState)
	statesRest.GET("/", getAllStates)
	statesRest.GET("/:name", getSingleState)
	statesRest.PUT("/:name", updateState)
	statesRest.DELETE("/:name", deleteState)
*/

//Communication has the functions to talk to PASSA server
type Communication struct {
	SchedulerURL string
}

//CreateState send the State object to PASSA
func (sc Communication) CreateState(state ymlparser.State) error {
	jsonValue, _ := json.Marshal(state)
	response, err := http.Post(sc.SchedulerURL+"/", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		return nil
	}
	return errors.New("Scheduler request failed")

}

//GetAllStates returns all states
func (sc Communication) GetAllStates() ([]ymlparser.State, error) {
	var returnedStates []ymlparser.State
	response, err := http.Get(sc.SchedulerURL + "/")
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 200 {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(data, &returnedStates)
		return returnedStates, nil
	}
	return nil, errors.New("Scheduler request failed")

}

//GetSingleState returns a single state
func (sc Communication) GetSingleState(stateName string) (ymlparser.State, error) {
	var returnedState ymlparser.State
	response, err := http.Get(sc.SchedulerURL + "/" + stateName)
	if err != nil {
		return ymlparser.State{}, err
	}
	if response.StatusCode == 200 {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return ymlparser.State{}, err
		}
		json.Unmarshal(data, &returnedState)
		return returnedState, nil
	}
	return ymlparser.State{}, errors.New("Scheduler request failed")

}

//UpdateState updates a old state to new server
func (sc Communication) UpdateState(oldStateName string, updateState ymlparser.State) error {
	jsonValue, _ := json.Marshal(updateState)
	response, err := http.Post(sc.SchedulerURL+"/"+oldStateName, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		return nil
	}
	return errors.New("Scheduler request failed")
}

//DeleteState deletes the state
func (sc Communication) DeleteState(deleteStateName string) error {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", sc.SchedulerURL+"/"+deleteStateName, nil)
	if err != nil {
		return err
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	return errors.New("Scheduler request failed")

}
