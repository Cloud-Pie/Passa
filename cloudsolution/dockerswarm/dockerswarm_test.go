package dockerswarm

import (
	"fmt"
	"reflect"
	"testing"

	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
)

const managerIP = "192.168.99.100"
const managerName = "myvm1"

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

func Test_changeState(t *testing.T) {

	//setup system

	//this test assumes two Vm's present so
	if len(listMachines()) != 2 {
		t.Skip("This test requires two machines to be present")
	}
	ds := NewSwarmManager("192.168.99.100")
	type args struct {
		wantedState ymlparser.Service
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "same state",
			args: args{
				wantedState: ymlparser.Service{Name: "graf_web", Scale: "2"},
			},
			want: []string{"myvm1", "myvm2"},
		}, {
			name: "add new machine",
			args: args{
				wantedState: ymlparser.Service{Name: "graf_web", Scale: "3"},
			},
			want: []string{"myvm1", "myvm2", "myvm3"},
		}, {
			name: "remove machine",
			args: args{
				wantedState: ymlparser.Service{Name: "graf_web", Scale: "2"},
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
