package fakecloud

import (
	"time"

	"github.com/Cloud-Pie/Passa/cloudsolution"
	"github.com/Cloud-Pie/Passa/ymlparser"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("passa")

type MockCloud struct {
	lastDeployedState ymlparser.State
}

func NewMockCloudManager(startingState ymlparser.State) MockCloud {
	m := MockCloud{}
	m.lastDeployedState = startingState
	log.Info("Mock Cloud created")
	return m
}

func (m MockCloud) GetLastDeployedState() ymlparser.State {
	return m.lastDeployedState
}
func (m MockCloud) ChangeState(wantedState ymlparser.State) cloudsolution.CloudManagerInterface {

	time.Sleep(5 * time.Second)
	m.lastDeployedState = wantedState
	return m
}

func (m MockCloud) GetActiveState() ymlparser.State {
	return m.lastDeployedState
}

func (m MockCloud) CheckState() bool {
	return true
}
