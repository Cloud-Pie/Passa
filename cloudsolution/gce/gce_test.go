package gce

import (
	"fmt"
	"testing"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

/*
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
*/
func Test_addLimits(t *testing.T) {
	wantedState := ymlparser.State{
		Services: ymlparser.Service{
			"nginx": ymlparser.ServiceInfo{Replicas: 2, CPU: "100m", Memory: "80Mi"},
		},
		VMs: ymlparser.VM{
			"t2.micro": 3,
			"t2.large": 3,
		},
	}
	for name, serviceInfo := range wantedState.Services {

		log.Info("Sending SCALE command for %s:%d", name, serviceInfo.Replicas)
		cmd := fmt.Sprintf(scaleContainersCommand, name, serviceInfo.Replicas, serviceInfo.CPU, serviceInfo.Memory)
		fmt.Println(cmd)
		//	exec.Command("sh", "-c", cmd).Output()
	}
}
