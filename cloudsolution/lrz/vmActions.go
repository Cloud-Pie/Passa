package lrz

import (
	"fmt"
	"io/ioutil"

	"os/exec"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Cloud-Pie/Passa/ymlparser"
	"k8s.io/client-go/kubernetes"
)

var types = []string{"t1.micro", "m1.nano", "m1.small", "m1.large", "m1.xlarge", "m2.xlarge", "c1.medium", "c1.xlarge", "m2.2xlarge", "m2.4xlarge", "cc1.4xlarge"}

const ec2URL = "https://www.cloud.mwn.de:22"
const privateKeyLine = "-----BEGIN RSA PRIVATE KEY-----"
const keyName = "passakey"
const keyFileName = "lrzkey.private"
const vmImage = "ami-00002826"

const createKeypairCommand = "econe-create-keypair %s -I %s -S %s -U %s"
const runInstanceCommand = "euca-run-instances -t %s -k %s -n %v -f %s %s -I %s -S %s -U %s"
const getInstancesCommand = "euca-describe-instances -I %s -S %s -U %s | grep running"
const terminateInstancesCommand = "euca-terminate-instances %s -I %s -S %s -U %s"

type econe struct {
	username   string
	password   string
	masterNode string
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
	_, err := exec.Command("sh", "-c", fmt.Sprintf(runInstanceCommand, templateType, keyName, vmNum, scriptFilename, vmImage, ec.username, ec.password, ec2URL)).Output()

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

func (ec econe) deleteMachine(currentVMState []string, templateType string, numToDelete int, kube *kubernetes.Clientset) error {

	var machineIDs []string
	var machineNames []string
	index := 0
	for _, line := range currentVMState {
		if strings.Contains(line, templateType) {
			mID := strings.Fields(line)[1]
			mName := strings.Split(strings.Fields(line)[3], ".")[0]

			if strings.Contains(line, ec.masterNode) {
				log.Error("%s is MASTER, can't delete", mName)
				log.Notice("MASTER node is running on %s", templateType)

			} else {
				machineNames = append(machineNames, mName)
				machineIDs = append(machineIDs, mID)
				index++
			}

			if index == numToDelete { //early exit
				break
			}
		}

	}
	f := strings.Join(machineIDs, " , ")
	c := fmt.Sprintf(terminateInstancesCommand, f, ec.username, ec.password, ec2URL)
	exec.Command("sh", "-c", c).Output()

	for _, machineName := range machineNames {
		log.Info("deleting machine %v", machineName)
		kube.CoreV1().Nodes().Delete(machineName, &metav1.DeleteOptions{})
	}

	return nil
}

func (ec econe) scaleVms(wantedVms []ymlparser.VM, kube *kubernetes.Clientset) {
	currentVms := ec.getVMs()
	wantedMap := make(map[string]int)
	currentMap := make(map[string]int)
	//wanted - current
	for _, vm := range wantedVms {
		wantedMap[vm.Type] = vm.Scale
	}

	for _, vm := range currentVms {
		currentMap[vm.Type] = vm.Scale
	}

	diffMap := make(map[string]int)
	for _, t := range types {
		if _, found := wantedMap[t]; found {
			diffMap[t] = wantedMap[t] - currentMap[t]
			log.Notice("changing state of %s\n", t)
		} else {
			log.Info("No change in %s\n", t)
		}
	}

	currentVMState, _ := exec.Command("sh", "-c", fmt.Sprintf(getInstancesCommand, ec.username, ec.password, ec2URL)).Output()

	a := strings.Split(string(currentVMState[:]), "\n")
	log.Debug("%v", diffMap)
	for changingTypes := range diffMap {
		numDiff := diffMap[changingTypes]
		switch {
		case numDiff == 0:
			//Do nothing
		case numDiff > 0:
			go ec.createNewVM(changingTypes, numDiff) //different type of VMs can be created in parallel
		case numDiff < 0:
			//delete machines
			//		deleteMachines(t, numDiff)
			go ec.deleteMachine(a, changingTypes, -numDiff, kube) //different type of VMs can be deleted in parallel
		}
	}

}
