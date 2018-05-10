package cloudsolution

import (
	"fmt"
	"reflect"
	"testing"

	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
)

const managerIP = "192.168.99.100"

func Test_getWorkerToken(t *testing.T) {

	fmt.Println(getWorkerToken(managerIP))
}

func Test_Integration(t *testing.T) {
	newMachineName := "myvm2"
	createNewMachine(newMachineName)
	newIP := getNewMachineIP(newMachineName)
	joinToken := getWorkerToken(managerIP)
	fmt.Println(addToSwarm(joinToken, newIP, managerIP, newMachineName))

}

func Test_createNewMachine(t *testing.T) {
	fmt.Printf("%s", createNewMachine("myvm2"))
}
func Test_deleteMachine(t *testing.T) {
	mn := "myvm2"
	fmt.Printf("%s", deleteMachine(mn))
}

func Test_listMachines(t *testing.T) {
	fmt.Printf("%s", listMachines())
	fmt.Printf("%v", len(listMachines()))

}

func Test_changeState(t *testing.T) {

	//setup system
	ds := NewSwarmManager("192.168.99.100")
	//this test assumes two Vm's present so
	if len(listMachines()) != 2 {
		t.Skip("This test requires two machines to be present")
	}
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
