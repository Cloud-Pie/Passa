//Package cloudsolution provides function for docker swarm
package cloudsolution

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"

	"golang.org/x/crypto/ssh"
)

//DockerSwarm keeps joinToken and managerIP of the system
type DockerSwarm struct {
	joinToken          string
	managerIP          string
	managerMachineName string
}

const machinePRefix = "myvm"

//NewSwarmManager returns a dockerswarm manager
func NewSwarmManager(managerIP string) DockerSwarm {

	managerName := "myvm1"
	return DockerSwarm{
		managerIP: managerIP,
		joinToken: getWorkerToken(managerIP, managerName), managerMachineName: managerName}
}

func createNewMachine(machineName string) []byte {
	cmd := exec.Command("sh", "-c", "docker-machine create --driver virtualbox "+machineName)

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

func getNewMachineIP(machineName string) string {
	newIP, err := exec.Command("sh", "-c", "docker-machine ip "+machineName).Output()

	if err != nil {
		panic(err)
	}
	return strings.Trim(string(newIP[:]), "\n")
}

func getWorkerToken(managerIP string, managerName string) string {

	session := getSSHSession(managerIP, managerName)

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("docker swarm join-token --quiet worker"); err != nil {
		log.Fatal("Failed to run:" + err.Error())
	}

	return b.String()
}

func (ds DockerSwarm) addToSwarm(newMachineIP string, machineName string) string {

	session := getSSHSession(newMachineIP, machineName)

	defer session.Close()

	var b bytes.Buffer

	session.Stdout = &b
	swarmCommand := fmt.Sprintf("docker swarm join --token %s %s:2377", strings.Trim(ds.joinToken, "\n"), ds.managerIP)
	fmt.Println(swarmCommand)
	if err := session.Run(swarmCommand); err != nil {
		log.Fatal("Failed to run:" + err.Error())
	}

	return b.String()
}

func (ds DockerSwarm) scaleContainers(containerName string, scaleNum string) string {

	session := getSSHSession(ds.managerIP, ds.managerMachineName)
	defer session.Close()

	//var b bytes.Buffer
	//session.Stdout = &b
	stdout, _ := session.StdoutPipe()

	scalingCommand := fmt.Sprintf("docker service scale %s=%s", containerName, scaleNum)

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

func deleteMachine(machineName string) []byte {
	cmd := exec.Command("sh", "-c", "docker-machine rm "+machineName+" -y")

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

func listMachines() []string {
	cmd := exec.Command("sh", "-c", "docker-machine ls -q")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	machinesList := strings.Split(strings.Trim(string(out[:]), "\n"), "\n")

	return machinesList
}

//ChangeState changes the state of the system
func (ds DockerSwarm) ChangeState(wantedState ymlparser.Service) []string {

	currentState := listMachines()
	scaleInt, err := strconv.Atoi(wantedState.Scale)
	if err != nil {
		panic(err)
	}
	difference := len(currentState) - scaleInt
	fmt.Println(difference)
	if difference == 0 { //keep the state as is
		return currentState
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
			newMachineName := fmt.Sprintf("%s%v", machinePRefix, len(currentState)+i+1)
			fmt.Println(newMachineName)
			go func() {
				defer wg.Done()
				createNewMachine(newMachineName)
				newIP := getNewMachineIP(newMachineName)

				ds.addToSwarm(newIP, newMachineName)

			}()

		}
		wg.Wait()
		ds.scaleContainers(wantedState.Name, wantedState.Scale)

	}

	return listMachines()
}

func (ds DockerSwarm) removeFromSwarm(machineName string) string {

	session := getSSHSession(ds.managerIP, ds.managerMachineName)
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("docker node rm -f " + machineName); err != nil {
		log.Fatal("Failed to run:" + err.Error())
	}

	return b.String()
}

func getSSHSession(machineIP string, machineName string) *ssh.Session {
	keyFile := fmt.Sprintf("%s/.docker/machine/machines/%s/id_rsa", os.Getenv("HOME"), machineName)
	key, err := ioutil.ReadFile(keyFile)

	signer, err := ssh.ParsePrivateKey(key)
	config := &ssh.ClientConfig{
		User: "docker",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //FIXME: fix security
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
