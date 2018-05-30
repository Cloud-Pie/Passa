package dockerswarm

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

const managerIP = "192.168.99.100"

func Test_getWorkerToken(t *testing.T) {
	if len(listMachines()) == 0 {
		t.Skip("No manager to get token")
	}
	fmt.Println(getWorkerToken(managerIP, managerName))
}

func Test_createNewMachine(t *testing.T) {
	if len(listMachines()) == 0 {
		t.Skip("No machine present")
	}
	fmt.Printf("%s", createNewMachine("myvm2"))
}
func Test_deleteMachine(t *testing.T) {
	if len(listMachines()) < 2 {
		t.Skip("No machine to delete")
	}
	mn := "myvm2"
	fmt.Printf("%s", deleteMachine(mn))
}

func Test_listMachines(t *testing.T) {
	machines := listMachines()
	fmt.Printf("%v", machines)

	fmt.Printf("%v", len(machines))

}

func TestChangeState(t *testing.T) {

	//setup system

	//this test assumes two Vm's present so
	if len(listMachines()) != 2 {
		t.Skip("This test requires two machines to be present")
	}
	ds := NewSwarmManager("192.168.99.100")
	type args struct {
		wantedState ymlparser.State
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "same state",
			args: args{
				wantedState: ymlparser.State{
					Services: []ymlparser.Service{{Name: "relax_web", Scale: 2}},
					VMs:      []ymlparser.VM{{Type: "medium", Scale: 2}},
				},
			},
			want: []string{"myvm1", "myvm2"},
		}, {
			name: "add new machine",
			args: args{
				wantedState: ymlparser.State{
					Services: []ymlparser.Service{{Name: "relax_web", Scale: 2}},
					VMs:      []ymlparser.VM{{Type: "medium", Scale: 3}},
				},
			},
			want: []string{"myvm1", "myvm2", "myvm3"},
		}, {
			name: "remove machine",
			args: args{
				wantedState: ymlparser.State{
					Services: []ymlparser.Service{{Name: "relax_web", Scale: 2}},
					VMs:      []ymlparser.VM{{Type: "medium", Scale: 2}},
				},
			},
			want: []string{"myvm1", "myvm2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ds.ChangeState(tt.args.wantedState); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangeState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerSwarm_GetActiveState(t *testing.T) {
	if len(listMachines()) == 0 {
		t.Skip("This test requires two machines to be present")
	}
	swarmManager := NewSwarmManager(managerIP)
	//swarmManager := DockerSwarm{}
	fmt.Println("swarmManager created")
	want := ymlparser.State{
		Services: []ymlparser.Service{{Name: "relax_web", Scale: 1}, {Name: "relax_visualizer", Scale: 1}},
		VMs:      []ymlparser.VM{{Type: machinePrefix, Scale: 1}},
	}

	sort.Slice(want.Services, func(i, j int) bool {
		return want.Services[i].Name > want.Services[j].Name
	})

	got := swarmManager.GetLastDeployedState()

	if !reflect.DeepEqual(got, want) {
		fmt.Printf("got: %#v \nwan: %#v\n", got, want)

		t.Fail()
	}

	//let's add a service
	wantedState := ymlparser.State{
		Services: []ymlparser.Service{{Name: "relax_web", Scale: 5}, {Name: "relax_visualizer", Scale: 1}},
		VMs:      []ymlparser.VM{{Type: machinePrefix, Scale: 1}},
	}

	changed := swarmManager.ChangeState(wantedState)
	if !reflect.DeepEqual(changed, wantedState) {
		fmt.Printf("changed\ngot: %#v \nwan: %#v\n", changed, wantedState)
		t.Fail()
	}

	wantedState = ymlparser.State{
		Services: []ymlparser.Service{{Name: "relax_web", Scale: 1}},
		VMs:      []ymlparser.VM{{Type: "medium", Scale: 1}},
	}
	swarmManager.ChangeState(wantedState) //back to normal
}
