package cloudsolution

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

const machineNAME = "myvm2"

//CreateNewMachine creates a new virtual machine in localhost with virtualbox
func CreateNewMachine(machineName string) []byte {
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

//GetNewMachineIP returns the IP of the newly created virtual machine
func GetNewMachineIP() string {
	newIP, err := exec.Command("sh", "-c", "docker-machine ip myvm2").Output()

	if err != nil {
		panic(err)
	}
	return string(newIP[:])
}

//GetWorkerToken returns the join token to add worker to the swarm
func GetWorkerToken(managerIP string) string {

	keyFile := fmt.Sprintf("%s/.docker/machine/machines/%s/id_rsa", os.Getenv("HOME"), "myvm1")
	key, err := ioutil.ReadFile(keyFile)

	signer, err := ssh.ParsePrivateKey(key)
	config := &ssh.ClientConfig{
		User: "docker",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //FIXME: fix security
	}

	fmt.Println(managerIP)
	client, err := ssh.Dial("tcp", managerIP+":22", config)

	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	session, err := client.NewSession()

	if err != nil {
		log.Fatal("Failed to session: ", err)
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("docker swarm join-token --quiet worker"); err != nil {
		log.Fatal("Failed to run:" + err.Error())
	}

	return b.String()
}

//AddToSwarm add newly created VM to docker swarm
func AddToSwarm(joinToken string, newMachineIP string, managerIP string) string {
	keyFile := fmt.Sprintf("%s/.docker/machine/machines/%s/id_rsa", os.Getenv("HOME"), machineNAME)
	fmt.Println(keyFile)
	key, err := ioutil.ReadFile(keyFile)

	signer, err := ssh.ParsePrivateKey(key)
	config := &ssh.ClientConfig{
		User: "docker",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //FIXME: fix security
	}

	fmt.Println(newMachineIP)
	client, err := ssh.Dial("tcp", newMachineIP+":22", config)

	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	session, err := client.NewSession()

	if err != nil {
		log.Fatal("Failed to session: ", err)
	}

	defer session.Close()

	var b bytes.Buffer

	session.Stdout = &b
	swarmCommand := fmt.Sprintf("docker swarm join --token %s %s:2377", strings.Trim(joinToken, "\n"), managerIP)
	fmt.Println(swarmCommand)
	if err := session.Run(swarmCommand); err != nil {
		log.Fatal("Failed to run:" + err.Error())
	}

	return b.String()
}

//ScaleContainers scales the containers to desired state
func ScaleContainers(managerIP string, containerName string, scaleNum string) string {
	keyFile := fmt.Sprintf("%s/.docker/machine/machines/%s/id_rsa", os.Getenv("HOME"), "myvm1")
	key, err := ioutil.ReadFile(keyFile)

	signer, err := ssh.ParsePrivateKey(key)
	config := &ssh.ClientConfig{
		User: "docker",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //FIXME: fix security
	}

	fmt.Println(managerIP)
	client, err := ssh.Dial("tcp", managerIP+":22", config)

	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	session, err := client.NewSession()

	if err != nil {
		log.Fatal("Failed to session: ", err)
	}

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

//DeleteMachine deletes the VM with the given name
func DeleteMachine(machineName string) []byte {
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

func changeState(wantedState int) []string {

	currentState := listMachines()
	difference := len(currentState) - wantedState
	fmt.Println(difference)
	if difference == 0 { //keep the state as is
		return currentState
	} else if difference > 0 { //lets delete some machines
		for i := 0; i < difference; i++ {
			lastCompName := currentState[len(currentState)-1]
			DeleteMachine(lastCompName)
			currentState = listMachines()
		}
	} else { //difference <0 , lets add some machines
		for i := 0; i < -1*difference; i++ {
			newMachineName := fmt.Sprintf("myvm%v", len(currentState)+1)
			CreateNewMachine(newMachineName)
			currentState = listMachines()
		}
	}

	return currentState
}
