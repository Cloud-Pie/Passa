package cloudsolution

import (
	"fmt"
	"reflect"
	"testing"
)

const managerIP = "192.168.99.100"

func Test_getWorkerToken(t *testing.T) {

	fmt.Println(GetWorkerToken(managerIP))
}

func Test_Integration(t *testing.T) {
	newMachineName := "myvm2"
	CreateNewMachine(newMachineName)
	newIP := GetNewMachineIP(newMachineName)
	joinToken := GetWorkerToken(managerIP)
	fmt.Println(AddToSwarm(joinToken, newIP, managerIP, newMachineName))

}

func Test_createNewMachine(t *testing.T) {
	fmt.Printf("%s", CreateNewMachine("myvm2"))
}
func Test_deleteMachine(t *testing.T) {
	mn := "myvm2"
	fmt.Printf("%s", DeleteMachine(mn))
}

func Test_listMachines(t *testing.T) {
	fmt.Printf("%s", listMachines())
	fmt.Printf("%v", len(listMachines()))

}

func Test_changeState(t *testing.T) {

	//this test assumes two Vm's present so
	if len(listMachines()) != 2 {
		t.Skip("This test requires two machines to be present")
	}
	type args struct {
		wantedState int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "same state",
			args: args{
				wantedState: 2,
			},
			want: []string{"myvm1", "myvm2"},
		}, {
			name: "add new machine",
			args: args{
				wantedState: 3,
			},
			want: []string{"myvm1", "myvm2", "myvm3"},
		}, {
			name: "remove machine",
			args: args{
				wantedState: 2,
			},
			want: []string{"myvm1", "myvm2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChangeState(tt.args.wantedState); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangeState() = %v, want %v", got, tt.want)
			}
		})
	}
}
