package lrz

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

func Test_econe_createNewKeypair(t *testing.T) {

}

func Test_VM_State(t *testing.T) {
	want := []ymlparser.VM{
		{
			Type:  types[0],
			Scale: 10,
		},

		{
			Type:  types[2],
			Scale: 2,
		},
	}

	current := []ymlparser.VM{
		{
			Type:  types[0],
			Scale: 1,
		},
		{
			Type:  types[1],
			Scale: 2,
		},
		{
			Type:  types[2],
			Scale: 5,
		},
	}

	diff := vmCalc(want, current)
	fmt.Printf("%v", diff)
	e := econe{}
	for _, t := range types {
		numDiff := diff[t]
		switch {
		case numDiff == 0:
			//do nothing
		case numDiff > 0:
			e.createNewVM(t, numDiff)
		case numDiff < 0:
			//delete machines
			//		deleteMachines(t, numDiff)
		}
	}
}

func vmCalc(wantedVms []ymlparser.VM, currentVms []ymlparser.VM) map[string]int {

	wantedMap := make(map[string]int)
	currentMap := make(map[string]int)
	//wanted - current
	for _, vm := range wantedVms {
		wantedMap[vm.Type] = vm.Scale
	}
	fmt.Printf("%v", wantedMap)
	for _, vm := range currentVms {
		currentMap[vm.Type] = vm.Scale
	}

	fmt.Printf("%v", currentMap)

	diffMap := make(map[string]int)
	for _, t := range types {
		diffMap[t] = wantedMap[t] - currentMap[t]
	}

	fmt.Printf("%v", diffMap)
	return diffMap
}

func Test_d_Machines(t *testing.T) {

	template := "m1.small"
	numToDelete := 0
	out := `INSTANCE	i-00039599			vm-10-155-209-58.cloud.mwn.de	running	none	39599		m1.large	2018-05-31T19:41:43+02:00	default	eki-EA801065	eri-1FEE1144		monitoring-disabled		10.155.209.58
	INSTANCE	i-00039931	ami-00002826		vm-10-155-209-45.cloud.mwn.de	running	passakey	39931		m1.nano	2018-06-19T16:05:15+02:00	default	eki-EA801065	eri-1FEE1144monitoring-disabled		10.155.209.45
	INSTANCE	i-00039932	ami-00002826		vm-10-155-209-61.cloud.mwn.de	running	passakey	39932		m1.small	2018-06-19T16:05:15+02:00	default	eki-EA801065	eri-1FEE1144		monitoring-disabled		10.155.209.61
	INSTANCE	i-00039933	ami-00002826		vm-10-155-209-49.cloud.mwn.de	running	passakey	39933		m1.nano	2018-06-19T16:05:15+02:00	default	eki-EA801065	eri-1FEE1144monitoring-disabled		10.155.209.49`

	a := strings.Split(string(out[:]), "\n")
	var machineIDs []string
	for _, line := range a {
		if strings.Contains(line, template) {
			mID := strings.Fields(line)[1] //1 because id is there
			mName := strings.Split(strings.Fields(line)[3], ".")[0]
			fmt.Println(mName)
			fmt.Println(mID)
			machineIDs = append(machineIDs, mID)
		}

	}
	f := strings.Join(machineIDs[:numToDelete], " , ")
	c := fmt.Sprintf("euca-terminate-instances %s -I di57dev -S e727d1464ae12436e899a726da5b2f11d8381b26 -U https://www.cloud.mwn.de:22", f)
	fmt.Printf("%s", c)

}
