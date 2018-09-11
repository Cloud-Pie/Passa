package gce

import (
	"fmt"
	"os/exec"

	"github.com/Cloud-Pie/Passa/cloudsolution"
	"github.com/Cloud-Pie/Passa/ymlparser"
	"github.com/op/go-logging"
)

//gcloud container clusters resize [CLUSTER_NAME] --node-pool [POOL_NAME] --size [SIZE]
const resizeClusterCommand = "gcloud container clusters resize %s --node-pool %s --size %d"

/*
* NOTE: gcloud info --format=json --filter="config.core"
*	same for kubectl
*
 */
const scaleContainersCommand = "kubectl scale deployment %s --replicas %d"

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
	if isCommandAvailable("gcloud") {
		log.Info("gcloud is active")
	} else {

		log.Critical("gcloud command not found, exiting...")

	}
	accountName, _ := exec.Command("sh", "-c", getAccount).Output()
	log.Debug("Authenticated as: %s", string(accountName))

	cs := GCE{
		clusterName: cn,
	}
	cs.lastDeployedState = cs.GetLastDeployedState()
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
	return GCE{}
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

func (g GCE) scaleVms(wantedVMs ymlparser.VM) { //TODO:
	for t, s := range wantedVMs {
		log.Info("Sending RESIZE command for %s:%d", t, s)
		exec.Command(fmt.Sprintf(resizeClusterCommand, g.clusterName, t, s))
	}
}

func (g GCE) scaleContainers(wantedContainers ymlparser.Service) {
	for name, replicaCount := range wantedContainers {
		log.Info("Sending SCALE command for %s:%d", name, replicaCount)
		exec.Command(fmt.Sprint(scaleContainersCommand, name, replicaCount)) //TODO: modify command for memory & cpu
	}
}

func (g GCE) getVMs() ymlparser.VM {
	out, _ := exec.Command(getDeploymentsCommand).Output()
	//TODO: parse out
	log.Info(string(out))
	return ymlparser.VM{}
}
func (g GCE) getServices() ymlparser.Service {
	out, _ := exec.Command(getNodesCommand).Output()
	//TODO: parse output
	log.Info(string(out))
	return ymlparser.Service{}
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
