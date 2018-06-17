package lrz

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

var types = []string{"m1.small", "m1.large", "m1.nano"}

const ec2URL = "https://www.cloud.mwn.de:22"
const privateKeyLine = "-----BEGIN RSA PRIVATE KEY-----"
const keyName = "passakey"
const keyFileName = "lrzkey.private"
const vmImage = "ami-00000001"

const createKeypairCommand = "econe-create-keypair %s -K %s -S %s -U %s" //FIXME: -K to -I for euca
const runInstanceCommand = "euca-run-instances -t %s -k %s -n %v %s -I %s -S %s -U %s"
const getInstancesCommand = "euca-describe-instances -I %s -S %s -U %s | grep running"
const terminateInstancesCommand = "euca-terminate-instances %s -I %s -S %s -U %s"

type econe struct {
	username string
	password string
}

func (ec econe) createNewKeypair() error {

	output, err := exec.Command("sh", "-c", fmt.Sprintf(createKeypairCommand, keyName, ec.username, ec.password, ec2URL)).Output()
	if err != nil {
		panic(err)
	}
	prKey := strings.Split(string(output), privateKeyLine)[1]
	d1 := []byte(privateKeyLine + prKey)
	err = ioutil.WriteFile(keyFileName, d1, 400)

	return err
}

func (ec econe) createNewVM(templateType string, vmNum int) error {
	_, err := exec.Command("sh", "-c", fmt.Sprintf(runInstanceCommand, templateType, keyName, vmNum, vmImage, ec.username, ec.password, ec2URL)).Output()

	return err
}

func (ec econe) getVMs() []ymlparser.VM {

	vms := []ymlparser.VM{}
	// only get running machines
	out, _ := exec.Command("sh", "-c", fmt.Sprintf(getInstancesCommand, ec.username, ec.password, ec2URL)).Output()

	for _, t := range types {

		vms = append(vms, ymlparser.VM{
			Type:  t,
			Scale: strings.Count(string(out), t),
		})
	}
	return vms
}

func (ec econe) deleteMachine(currentVMState []string, templateType string, numToDelete int) error {

	var machineIDs []string
	for _, line := range currentVMState {
		if strings.Contains(line, templateType) {
			mID := strings.Fields(line)[1]
			fmt.Println(mID)
			machineIDs = append(machineIDs, mID)
		}

	}
	f := strings.Join(machineIDs[:numToDelete], " , ")
	c := fmt.Sprintf(terminateInstancesCommand, f, ec.username, ec.password, ec2URL)
	fmt.Printf("%s", c)

	return nil
}

func (ec econe) scaleVms(wantedVms []ymlparser.VM, currentVms []ymlparser.VM) {
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

	currentVMState, _ := exec.Command("sh", "-c", fmt.Sprintf(getInstancesCommand, ec.username, ec.password, ec2URL)).Output()

	a := strings.Split(string(currentVMState[:]), "\n")
	fmt.Printf("%v", diffMap)
	for _, t := range types {
		numDiff := diffMap[t]
		switch {
		case numDiff == 0:
			//Do nothing
		case numDiff > 0:
			ec.createNewVM(t, numDiff)
		case numDiff < 0:
			//delete machines
			//		deleteMachines(t, numDiff)
			ec.deleteMachine(a, t, -numDiff)
		}
	}
}
