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
	want := ymlparser.VM{
		"m1.large": 10,
		"m1.nano":  2,
	}

	current := ymlparser.VM{
		"m1.large": 1,
		"m1.small": 2,
		"m1.nano":  5,
	}

	diff := vmCalc(want, current)

	for k, v := range diff {
		fmt.Printf("%v -> %v\n", k, v)

	}

}

func vmCalc(wantedMap ymlparser.VM, currentMap ymlparser.VM) ymlparser.VM {

	//wanted - current

	diffMap := make(map[string]int)
	for k := range wantedMap {
		if v, ok := currentMap[k]; ok {
			diffMap[k] = wantedMap[k] - v
		}
	}

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
