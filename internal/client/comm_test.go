package client

/*
import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/Cloud-Pie/Passa/database"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

func TestCommunication(t *testing.T) {

	r := SetupServer(testManager{}, make(chan *ymlparser.State, 30))

	go r.Run()

	log.Println("Server Set up")
	log.Println("testing... Started")

	comm := Communication{
		SchedulerURL: "http://localhost:8080/api/states",
	}
	stateToTest := ymlparser.State{
		ISODate:  time.Now(),
		Name:     "test-State",
		Services: append([]ymlparser.Service{}, ymlparser.Service{Name: "test-service", Scale: 10}),
	}
	//Create
	err := comm.CreateState(stateToTest)

	if err != nil {
		log.Println("Create Failed")
		t.Fail()
	}
	//GetAll
	returnedStates, err := comm.GetAllStates()
	if err != nil || !reflect.DeepEqual(returnedStates, database.ReadAllStates()) {
		log.Printf("GetAll Failed\n %+v \n %+v", returnedStates, database.ReadAllStates())
		t.Fail()
	}
	//GetSingle
	returnedState, err := comm.GetSingleState("test-State")
	if !reflect.DeepEqual(returnedState.Services, stateToTest.Services) {

		log.Printf("GetSingle Failed. \n%+v \n%+v", returnedState.ISODate, stateToTest.ISODate)
		t.Fail()
	}

	//UpdateState
	stateToUpdate := stateToTest
	stateToUpdate.ISODate = time.Now()
	err = comm.UpdateState(stateToTest.Name, stateToUpdate)
	if err != nil {
		log.Println("UpdateState Failed")
		t.Fail()
	}

	//
	err = comm.DeleteState(stateToTest.Name)
	if err != nil {
		log.Println("DeleteState Failed")
		t.Fail()
	}
	log.Println("testing... End")
}
*/
