package gce

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Cloud-Pie/Passa/cloudsolution"
	"github.com/Cloud-Pie/Passa/ymlparser"
	"github.com/op/go-logging"
)

//gcloud container clusters resize [CLUSTER_NAME] --node-pool [POOL_NAME] --size [SIZE]
const resizeClusterCommand = "gcloud container clusters  resize %s --node-pool %s --size %d -q"

var types = []string{"t2.micro", "t2.large"}

const scaleContainersCommand = "kubectl scale deployment %s --replicas %d "

const getNodesCommand = "kubectl get nodes"

const getDeploymentsCommand = "kubectl get deployments"

const getAccount = "gcloud info --format='value(config.account)'"

var log = logging.MustGetLogger("passa")

//GCE keeps configuration for Google Cloud Engine
type GCE struct {
	lastDeployedState ymlparser.State
	clusterName       string
}

//NewGCEManager return a new manager for GCE.
func NewGCEManager(cn string) GCE {
	if isCommandAvailable("gcloud") && isCommandAvailable("kubectl") {
		log.Info("Commands found: gcloud, kubectl")
	} else {

		log.Critical("gcloud or kubectl not found")

	}
	accountName, _ := exec.Command("sh", "-c", getAccount).Output()
	log.Debug("Authenticated as: %s", string(accountName))

	cs := GCE{
		clusterName: cn,
	}
	cs.lastDeployedState = cs.GetActiveState()
	log.Info("GCE manager created")
	return cs
}

//GetLastDeployedState returns last deployed state by the system
func (g GCE) GetLastDeployedState() ymlparser.State {
	return g.lastDeployedState
}

//ChangeState changes the state.
func (g GCE) ChangeState(wantedState ymlparser.State) cloudsolution.CloudManagerInterface {
	if wantedState.VMs != nil {
		g.scaleVms(wantedState.VMs)
	}
	g.scaleContainers(wantedState.Services)

	g.lastDeployedState = wantedState
	return g
}

//GetActiveState gets current state.
func (g GCE) GetActiveState() ymlparser.State {
	return ymlparser.State{
		VMs:      g.getVMs(),
		Services: g.getServices(),
	}
}

//CheckState checks the state of the cloud.
func (g GCE) CheckState() bool {
	weDeployed := g.GetLastDeployedState()
	real := g.GetActiveState()

	if areVMsCorrect(weDeployed.VMs, real.VMs) && areServicesCorrect(weDeployed.Services, real.Services) {

		return true
	}
	log.Error("ERROR:\ndepl: %#v\nreal: %#v", weDeployed, real)
	return false
}

func (g GCE) scaleVms(wantedVMs ymlparser.VM) {

	for t, s := range wantedVMs {
		t = strings.Replace(t, ".", "-", -1) // For AWS compatibility
		log.Info("Sending RESIZE command for %s:%d", t, s)
		cmd := fmt.Sprintf(resizeClusterCommand, g.clusterName, t, s)
		fmt.Println(cmd)
		exec.Command("sh", "-c", cmd).Output()
	}
}

func (g GCE) scaleContainers(wantedContainers ymlparser.Service) {

	for name, serviceInfo := range wantedContainers {
		log.Info("Sending SCALE command for %s:%d", name, serviceInfo.Replicas)
		cmd := fmt.Sprintf(scaleContainersCommand, name, serviceInfo.Replicas)

		fmt.Println(cmd)
		exec.Command("sh", "-c", cmd).Output()
	}
}

func (g GCE) getVMs() ymlparser.VM {
	vms := ymlparser.VM{}
	out, _ := exec.Command("sh", "-c", getNodesCommand).Output()

	for _, t := range types {
		searchString := strings.Replace(t, ".", "-", -1)
		vms[t] = strings.Count(string(out), searchString)

	}
	fmt.Println(vms)
	return vms

}

func (g GCE) getServices() ymlparser.Service {
	serviceMap := ymlparser.Service{}
	services, _ := exec.Command("sh", "-c", "kubectl get deployments").Output()

	a := strings.Split(string(services[:]), "\n")

	for _, line := range a[1 : len(a)-1] {
		serviceName := strings.Fields(line)[0]
		replicaCount := strings.Fields(line)[4]

		replicaCountInt, err := strconv.Atoi(replicaCount)
		if err != nil {
			panic(err)
		}
		serviceMap[serviceName] = ymlparser.ServiceInfo{Replicas: replicaCountInt, CPU: "", Memory: 0}
		log.Info("%v", serviceMap)

	}
	return serviceMap
}

func areVMsCorrect(deployedVMMap ymlparser.VM, realVMMap ymlparser.VM) bool {

	for key := range deployedVMMap {
		if deployedVMMap[key] != realVMMap[key] {
			return false
		}
	}
	return true
}

func areServicesCorrect(deployedServicesMap ymlparser.Service, realServicesMap ymlparser.Service) bool {

	for key := range deployedServicesMap {
		if deployedServicesMap[key] != realServicesMap[key] {
			return false
		}
	}
	return true
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
