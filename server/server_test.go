package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
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
	req, _ = http.NewRequest("GET", "/ui/states", nil)

	r.ServeHTTP(w, req)
	var returnedData *ymlparser.Config
	json.Unmarshal(w.Body.Bytes(), &returnedData)

	if !reflect.DeepEqual(c, returnedData) {
		t.Errorf("want: %v, got: %v", c, returnedData)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/ui/timeline", nil)

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fail()
	}

}

func TestServer(t *testing.T) {
	c := ymlparser.ParseStatesfile("../test/passa-states-test.yml")
	r := SetupServer(c)
	r.Run(":2000")
}
