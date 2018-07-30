//Package dockerswarm provides function for docker swarm
package dockerswarm

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/op/go-logging"

	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/Cloud-Pie/Passa/cloudsolution"

	"github.com/Cloud-Pie/Passa/ymlparser"
	"golang.org/x/crypto/ssh"
)

var log = logging.MustGetLogger("passa")

//DockerSwarm keeps joinToken and managerIP of the system
type DockerSwarm struct {
	joinToken           string
	managerIP           string
	managerMachineName  string
	lastDeployedState   ymlparser.State
	isActivelyDeploying bool
}

//machinePrefix makes sure all our machines have names like myvm1, myvm2, myvm3.
const (
	machinePrefix           = "myvm"
	managerName             = "myvm1"
	createNewMachineCommand = "docker-machine create --driver virtualbox %s"
	deleteMachineCommand    = "docker-machine rm %s -y"
	getIPCommand            = "docker-machine ip %s"
	getWorkerTokenCommand   = "docker swarm join-token --quiet worker"
	joinWorkerCommand       = "docker swarm join --token %s %s:2377"
	scaleServiceCommand     = "docker service scale %s=%v"
	listMachineCommand      = "docker-machine ls -q"
	removeFromSwarmCommand  = "docker node rm -f %s"
	dockerKeyLocation       = "%s/.docker/machine/machines/%s/id_rsa"
	getServiceCommand       = "docker service ls --format '{{.Name}} {{.Replicas}}'"
)

//NewSwarmManager returns a dockerswarm manager
func NewSwarmManager(managerIP string) cloudsolution.CloudManagerInterface {
	dc := DockerSwarm{
		managerIP:          managerIP,
		joinToken:          getWorkerToken(managerIP, managerName),
		managerMachineName: managerName,
	}
	dc.lastDeployedState = dc.GetActiveState()
	return dc
}

//CreateNewMachine creates new machine with the docker-machine command.
func createNewMachine(machineName string) []byte {
	cmd := exec.Command("sh", "-c", fmt.Sprintf(createNewMachineCommand, machineName))

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()

	out, _ := cmd.Output()
	return out
}

//getNewMachineIP returns the IP of the asked machine.
func getNewMachineIP(machineName string) string {
	newIP, err := exec.Command("sh", "-c", fmt.Sprintf(getIPCommand, machineName)).Output()

	if err != nil {
		panic(err)
	}
	return strings.Trim(string(newIP[:]), "\n")
}

//getWorkerToken returns the worker token required to join the swarm
func getWorkerToken(managerIP string, managerName string) string {

	session := getSSHSession(managerIP, managerName)

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(getWorkerTokenCommand); err != nil {
		log.Fatal("Failed to run:" + err.Error())
	}

	return b.String()
}

//addToSwarm add newly created vm to docker swarm
func (ds DockerSwarm) addToSwarm(newMachineIP string, machineName string) string {

	session := getSSHSession(newMachineIP, machineName)

	defer session.Close()

	var b bytes.Buffer

	session.Stdout = &b
	swarmCommand := fmt.Sprintf(joinWorkerCommand, strings.Trim(ds.joinToken, "\n"), ds.managerIP)
	fmt.Println(swarmCommand)
	if err := session.Run(swarmCommand); err != nil {
		log.Fatal("Failed to run:" + err.Error())
	}

	return b.String()
}

//scaleContainers give command to manager to scale the services.
func (ds DockerSwarm) scaleContainers(serviceName string, scaleNum int) string {

	session := getSSHSession(ds.managerIP, ds.managerMachineName)
	defer session.Close()

	//var b bytes.Buffer
	//session.Stdout = &b
	stdout, _ := session.StdoutPipe()

	scalingCommand := fmt.Sprintf(scaleServiceCommand, serviceName, scaleNum)

	if err := session.Start(scalingCommand); err != nil {
		log.Fatal("Failed to run:" + err.Error())
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	session.Wait()
	return ""
}

//deleteMachine deletes the machine.
func deleteMachine(machineName string) []byte {
	cmd := exec.Command("sh", "-c", fmt.Sprintf(deleteMachineCommand, machineName))

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()

	out, _ := cmd.Output()
	return out
}

//listMachines lists currently created machines.
//It assumes that all the machines created are running.
func listMachines() []string {
	cmd := exec.Command("sh", "-c", listMachineCommand)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	if len(out) == 0 {
		return []string{}
	}

	machinesList := strings.Split(strings.Trim(string(out[:]), "\n"), "\n")

	return machinesList
}

//removeFromSwarm removes the deleted vm from swarm
func (ds DockerSwarm) removeFromSwarm(machineName string) string {

	session := getSSHSession(ds.managerIP, ds.managerMachineName)
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(fmt.Sprintf(removeFromSwarmCommand, machineName)); err != nil {
		log.Fatal("Failed to run:" + err.Error())
	}

	return b.String()
}

//getSSHSession returns the SSH session to the machine with machineName and MachineIP.
func getSSHSession(machineIP string, machineName string) *ssh.Session {
	keyFile := fmt.Sprintf(dockerKeyLocation, os.Getenv("HOME"), machineName)
	key, err := ioutil.ReadFile(keyFile)

	signer, err := ssh.ParsePrivateKey(key)
	config := &ssh.ClientConfig{
		User: "docker",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", machineIP+":22", config)

	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	session, err := client.NewSession()

	if err != nil {
		log.Fatal("Failed to session: ", err)
	}

	return session
}

//ChangeState changes the state of the system
func (ds DockerSwarm) ChangeState(wantedState ymlparser.State) cloudsolution.CloudManagerInterface {
	ds.isActivelyDeploying = true
	if wantedState.VMs != nil { //There is no typing in docker swarm so take it like this.

		totalVM := 0
		for idx := range wantedState.VMs {
			totalVM += wantedState.VMs[idx]
		}
		//Scale machines
		currentState := listMachines()
		difference := len(currentState) - totalVM
		fmt.Println(difference)
		if difference == 0 { //keep the state as is
			fmt.Println("No new machine")
		} else if difference > 0 { //lets delete some machines
			var wg sync.WaitGroup
			wg.Add(difference)
			for i := 0; i < difference; i++ {
				lastCompName := currentState[len(currentState)-1-i]
				go func() {
					defer wg.Done()
					deleteMachine(lastCompName)
					ds.removeFromSwarm(lastCompName)
				}()

			}
			wg.Wait()
		} else { //difference <0 , lets add some machines
			var wg sync.WaitGroup
			wg.Add(-difference)
			for i := 0; i < -1*difference; i++ {
				newMachineName := fmt.Sprintf("%s%v", machinePrefix, len(currentState)+i+1)
				fmt.Println(newMachineName)
				go func() {
					defer wg.Done()
					createNewMachine(newMachineName)
					newIP := getNewMachineIP(newMachineName)

					ds.addToSwarm(newIP, newMachineName)

				}()

			}
			wg.Wait()
			//Scale containers
			fmt.Printf("%#v", wantedState)

		}
	} else {
		log.Debug("%s has no VM state, keeping current...", wantedState.Name)
	}
	for key := range wantedState.Services {
		ds.scaleContainers(key, wantedState.Services[key].Replicas)

	}

	ds.lastDeployedState = ds.GetActiveState()
	ds.isActivelyDeploying = false
	return ds
}

//GetActiveState gets the current state of the cloud
func (ds DockerSwarm) GetActiveState() ymlparser.State {

	return ymlparser.State{
		VMs:      ds.getMachines(),
		Services: ds.getServiceCount(),
	}
}

func (ds DockerSwarm) getServiceCount() ymlparser.Service {
	//return []ymlparser.Service{}

	session := getSSHSession(ds.managerIP, ds.managerMachineName)
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(getServiceCommand); err != nil {
		log.Fatal("Failed to run:" + err.Error())
	}

	servicesList := strings.Split(strings.Trim(string(b.String()[:]), "\n"), "\n")

	currentServices := ymlparser.Service{}
	for _, serviceString := range servicesList {
		serviceSplit := strings.Split(serviceString, " ")

		serviceCount, _ := strconv.Atoi(strings.Split(serviceSplit[1], "/")[0])

		currentServices[serviceSplit[0]] = ymlparser.ServiceInfo{Replicas: serviceCount}

	}

	return currentServices
}

//GetLastDeployedState returns the State that we believe is currently running in cloud
func (ds DockerSwarm) GetLastDeployedState() ymlparser.State {
	return ds.lastDeployedState
}

//CheckState compares the actual state and the state we have deployed.
func (ds DockerSwarm) CheckState() bool {

	weDeployed := ds.GetLastDeployedState()
	real := ds.GetActiveState() //SORT

	real.ISODate = weDeployed.ISODate //server return zero ISODate and is equal check fails otherwise
	if reflect.DeepEqual(weDeployed, real) {
		log.Info("State holds")
		return true
	}

	log.Error("ERROR: \ndepl: %#v\nreal: %#v\n", weDeployed, real)
	return false
}

func (ds DockerSwarm) getMachines() ymlparser.VM {

	return ymlparser.VM{
		machinePrefix: len(listMachines()),
	}

}
