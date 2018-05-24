package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

func TestSetupServer(t *testing.T) {
	c := ymlparser.ParseStatesfile("../test/passa-states-test.yml")
	r := SetupServer(c)
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
	c := ymlparser.ParseStatesfile("../test/passa-states-test.yml")
	r := SetupServer(c)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/states/", nil)
	r.ServeHTTP(w, req)

	var returnedData []ymlparser.State
	json.Unmarshal(w.Body.Bytes(), &returnedData)
	if !reflect.DeepEqual(c.States, returnedData) {
		t.Errorf("want: %v\ngot: %v\n", c.States, returnedData)
	}
}

func Test_createState(t *testing.T) {

	c := ymlparser.ParseStatesfile("../test/passa-states-test.yml")
	r := SetupServer(c)
	stateNum := len(c.States)
	tests := []struct {
		name          string
		stateToUpdate ymlparser.State
		returnedCode  int
	}{

		{
			name: "Valid State",
			stateToUpdate: ymlparser.State{
				Time:     "18-08-2018, 20:00:00 CEST",
				Name:     "test-State",
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
				Name: "Invalid State",
				Time: "18-08-2018, 20:00:00 CEST",
			},
			returnedCode: 422,
		},
		{
			name: "Valid State with Invalid Time",
			stateToUpdate: ymlparser.State{
				Time:     "dummy time",
				Name:     "test-State",
				Services: append([]ymlparser.Service{}, ymlparser.Service{Name: "test-service", Scale: 10}),
			},
			returnedCode: 422,
		},
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
				var returnedData []ymlparser.State
				json.Unmarshal(w.Body.Bytes(), &returnedData)
				if stateNum+1 != len(returnedData) {
					t.Errorf("want: %v\ngot: %v\n", stateNum+1, len(returnedData))
				}

			}
		})
	}

}

func Test_getSingleState(t *testing.T) {
	c := ymlparser.ParseStatesfile("../test/passa-states-test.yml")
	r := SetupServer(c)

	tests := []struct {
		name         string
		stateName    string
		returnedCode int
	}{

		{
			name:         "Existing State",
			stateName:    "state-with-7",
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
	c := ymlparser.ParseStatesfile("../test/passa-states-test.yml")
	r := SetupServer(c)

	tests := []struct {
		name          string
		stateToUpdate ymlparser.State
		stateName     string
		returnedCode  int
	}{

		{
			name: "Existing State To Update",
			stateToUpdate: ymlparser.State{
				Time: "18-08-2019, 15:45:33 CEST",
				Name: "update State",
			},
			returnedCode: 200,
			stateName:    "state-with-7",
		},
		{
			name: "Non Existing State To Update",
			stateToUpdate: ymlparser.State{
				Time: "18-08-2019, 15:45:33 CEST",
				Name: "Non existent State",
			},
			returnedCode: 422,
			stateName:    "non-existent",
		},
		{name: "Existing State To Update with Invalid Time",
			stateToUpdate: ymlparser.State{
				Time: "dummy time",
				Name: "state-with-dummy-time",
			},
			returnedCode: 422,
			stateName:    "state-with-6",
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
			req, _ := http.NewRequest("PUT", stateLink, bytes.NewBuffer(updateState))
			r.ServeHTTP(w, req)
			if w.Code != tt.returnedCode {
				t.Errorf("want: %v\ngot: %v\n", tt.returnedCode, w.Code)
			}
		})
	}
}

func Test_deleteState(t *testing.T) {
	c := ymlparser.ParseStatesfile("../test/passa-states-test.yml")
	r := SetupServer(c)

	tests := []struct {
		name           string
		stateName      string
		returnedLength int
		returnedCode   int
	}{

		{
			name:           "Existing State",
			stateName:      "state-with-7",
			returnedLength: len(c.States) - 1,
			returnedCode:   200,
		},
		{
			name:           "Non Existing State",
			stateName:      "non-state",
			returnedLength: len(c.States),
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
				var returnedData []ymlparser.State
				json.Unmarshal(w.Body.Bytes(), &returnedData)
				if tt.returnedLength != len(returnedData) {
					t.Errorf("want: %v\ngot: %v\n", tt.returnedLength, len(returnedData))
				}
			}
		})
	}
}

func Test_test(t *testing.T) {
	c := ymlparser.ParseStatesfile("../test/passa-states-test.yml")
	r := SetupServer(c)
	w := httptest.NewRecorder()
	myState := ymlparser.State{
		Time: "18-08-1994, 20:00:00 CEST",
		Name: "myState",
	}
	jsonState, err := json.Marshal(myState)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", jsonState)
	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonState))
	r.ServeHTTP(w, req)
}

func TestRoutes(t *testing.T) {
	c := ymlparser.ParseStatesfile("../test/passa-states-test.yml")
	r := SetupServer(c)
	fmt.Printf("%v", r.Routes())
}
