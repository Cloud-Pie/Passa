package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Cloud-Pie/Passa/cloudsolution"
	"github.com/Cloud-Pie/Passa/database"
	"github.com/Cloud-Pie/Passa/ymlparser"
)

func TestSetupServer(t *testing.T) {

	r := SetupServer(testManager{}, make(chan *ymlparser.State, 30))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fail()
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/ui/timeline", nil)

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fail()
	}

}
func Test_getAllStates(t *testing.T) {

	r := SetupServer(testManager{}, make(chan *ymlparser.State, 30))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/states/", nil)
	r.ServeHTTP(w, req)

	var returnedData []ymlparser.State
	json.Unmarshal(w.Body.Bytes(), &returnedData)
	if !reflect.DeepEqual(database.ReadAllStates(), returnedData) {
		t.Errorf("want: %v\ngot: %v\n", database.ReadAllStates(), returnedData)
	}
}

func Test_createState(t *testing.T) {

	tc := make(chan *ymlparser.State, 30)
	r := SetupServer(testManager{}, tc)
	stateNum := len(database.ReadAllStates())
	tests := []struct {
		name          string
		stateToUpdate ymlparser.State
		returnedCode  int
	}{

		{
			name: "Valid State",
			stateToUpdate: ymlparser.State{
				ISODate:  time.Now(),
				Name:     "state-with-6",
				Services: append([]ymlparser.Service{}, ymlparser.Service{Name: "test-service", Scale: 10}),
			},
			returnedCode: 200,
		},
		{
			name: "Invalid State without Time",
			stateToUpdate: ymlparser.State{
				Name: "Invalid State",
			},
			returnedCode: 422,
		},
		{
			name: "Invalid State without Service",
			stateToUpdate: ymlparser.State{
				Name:    "Invalid State",
				ISODate: time.Now(),
			},
			returnedCode: 422,
		}, //invalid string time make ISODate zero
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonState, err := json.Marshal(tt.stateToUpdate)
			if err != nil {
				panic(err)
			}
			req, _ := http.NewRequest("POST", "/api/states/", bytes.NewBuffer(jsonState))
			r.ServeHTTP(w, req)
			if w.Code != tt.returnedCode {
				t.Errorf("want: %v\ngot: %v\n", tt.returnedCode, w.Code)
			}

			if w.Code == http.StatusOK {
				w = httptest.NewRecorder()
				req, _ = http.NewRequest("GET", "/api/states/", nil)
				r.ServeHTTP(w, req)
				var returnedData []ymlparser.State
				json.Unmarshal(w.Body.Bytes(), &returnedData)
				if stateNum+1 != len(returnedData) {
					t.Errorf("want: %v\ngot: %v\n", stateNum+1, len(returnedData))
				}

				expected := <-tc
				if reflect.DeepEqual(expected, tt.stateToUpdate) {
					fmt.Printf("got: %#v\nwan: %#v", expected, tt.stateToUpdate)

				}
			}
		})
	}
}

func Test_getSingleState(t *testing.T) {

	r := SetupServer(testManager{}, make(chan *ymlparser.State, 30))

	tests := []struct {
		name         string
		stateName    string
		returnedCode int
	}{

		{
			name:         "Existing State",
			stateName:    "state-with-6",
			returnedCode: 200,
		},
		{
			name:         "Non Existing State",
			stateName:    "non-state",
			returnedCode: 422,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			stateLink := fmt.Sprintf("/api/states/%s", tt.stateName)
			req, _ := http.NewRequest("GET", stateLink, nil)
			r.ServeHTTP(w, req)
			if w.Code != tt.returnedCode {
				t.Errorf("want: %v\ngot: %v\n", tt.returnedCode, w.Code)
			}
		})
	}
}

func Test_updateState(t *testing.T) {

	r := SetupServer(testManager{}, make(chan *ymlparser.State, 30))

	tests := []struct {
		name          string
		stateToUpdate ymlparser.State
		stateName     string
		returnedCode  int
	}{

		{
			name: "Existing State To Update",
			stateToUpdate: ymlparser.State{
				ISODate: time.Now(),
				Name:    "update State",
			},
			returnedCode: 200,
			stateName:    "state-with-6",
		},
		{
			name: "Non Existing State To Update",
			stateToUpdate: ymlparser.State{
				ISODate: time.Now(),
				Name:    "Non existent State",
			},
			returnedCode: 422,
			stateName:    "non-existent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			stateLink := fmt.Sprintf("/api/states/%s", tt.stateName)
			updateState, err := json.Marshal(tt.stateToUpdate)
			if err != nil {
				panic(err)
			}
			req, _ := http.NewRequest("POST", stateLink, bytes.NewBuffer(updateState))
			r.ServeHTTP(w, req)
			if w.Code != tt.returnedCode {
				t.Errorf("want: %v\ngot: %v\n", tt.returnedCode, w.Code)
			}
		})
	}
}

func Test_deleteState(t *testing.T) {

	r := SetupServer(testManager{}, make(chan *ymlparser.State, 30))

	tests := []struct {
		name           string
		stateName      string
		returnedLength int
		returnedCode   int
	}{

		{
			name:           "Existing State",
			stateName:      "state-with-6",
			returnedLength: len(database.ReadAllStates()) - 1,
			returnedCode:   200,
		},
		{
			name:           "Non Existing State",
			stateName:      "non-state",
			returnedLength: len(database.ReadAllStates()),
			returnedCode:   422,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			stateLink := fmt.Sprintf("/api/states/%s", tt.stateName)
			req, _ := http.NewRequest("DELETE", stateLink, nil)
			r.ServeHTTP(w, req)
			if w.Code != tt.returnedCode {
				t.Errorf("want: %v\ngot: %v\n", tt.returnedCode, w.Code)
			}

			if w.Code == 200 {
				w = httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/api/states/", nil)
				r.ServeHTTP(w, req)

				var returnedData []ymlparser.State
				json.Unmarshal(w.Body.Bytes(), &returnedData)
				if len(returnedData) != tt.returnedLength {
					t.Fail()
				}
			}
		})
	}
}

func Test_test(t *testing.T) {

	r := SetupServer(testManager{}, make(chan *ymlparser.State, 30))
	w := httptest.NewRecorder()
	myTime := time.Now()
	myState := ymlparser.State{
		ISODate: myTime,
		Name:    "myTestState",
	}
	jsonState, err := json.Marshal(myState)
	if err != nil {
		panic(err)
	}
	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonState))
	r.ServeHTTP(w, req)

}

type testManager struct {
	name string
}

func (testManager) ChangeState(ymlparser.State) cloudsolution.CloudManagerInterface {

	return testManager{}
}
func (testManager) GetActiveState() ymlparser.State {

	return ymlparser.State{
		Name: "My Active State",
	}
}
func (testManager) GetLastDeployedState() ymlparser.State {
	return ymlparser.State{
		Name: "My Last Deployed State",
	}
}
func (testManager) CheckState() bool {

	return true
}

func Test_getCurrentState(t *testing.T) {
	tm := testManager{
		name: "test",
	}
	r := SetupServer(tm, make(chan *ymlparser.State, 30))
	r.Run(":7000")
}
