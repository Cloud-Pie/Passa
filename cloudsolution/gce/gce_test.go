package gce

import (
	"testing"

	"github.com/Cloud-Pie/Passa/cloudsolution"
	"github.com/Cloud-Pie/Passa/ymlparser"
)

func Test_NewGCEManager(t *testing.T) {
	wantedState := ymlparser.State{
		Services: ymlparser.Service{
			"nginx": ymlparser.ServiceInfo{Replicas: 2},
		},
		VMs: ymlparser.VM{
			"t2.micro": 3,
			"t2.large": 3,
		},
	}

	var cs cloudsolution.CloudManagerInterface
	cs = NewGCEManager("passa-cluster")
	cs = cs.ChangeState(wantedState)

	if !cs.CheckState() {
		t.Fail()
	}

}
