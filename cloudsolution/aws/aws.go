package aws

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/op/go-logging"

	"github.com/Cloud-Pie/Passa/cloudsolution"
	"github.com/Cloud-Pie/Passa/ymlparser"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

const scriptFilename = "awsscript.sh"
const bashCommand = "#!/bin/bash"
const deploymentTimeout = 120 * time.Second

var log = logging.MustGetLogger("passa")

//Lrz keeps the data needed for econe and kubernetes interfaces.
type Lrz struct {
	lastDeployedState ymlparser.State
	econe             econe
	kube              *kubernetes.Clientset
}

//NewAWSManager return a new manager for lrz.
func NewAWSManager(username, password, configFile string, joinCommand string) Lrz {

	config, err := clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		log.Fatal("Couldn't build config from file")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Cannot connect to kubernetes cluster , exiting...")
	}
	cs := Lrz{
		econe: econe{
			username: username,
			password: password,
		},
		kube: clientset,
	}
	log.Debug("Adding join token to file")
	data := []byte(fmt.Sprintf("%s\n%s", bashCommand, joinCommand))
	ioutil.WriteFile(scriptFilename, data, 0644)

	nodesList, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})

	for _, node := range nodesList.Items {
		for k := range node.GetLabels() {
			if strings.Contains(k, "master") {
				log.Warning("%s is MASTER", node.Name)
				cs.econe.masterNode = node.Name
			}
		}
	}
	cs.lastDeployedState = cs.GetActiveState()
	return cs
}

//ChangeState deploys wanted state to LRZ and kubernetes
func (l Lrz) ChangeState(wantedState ymlparser.State) cloudsolution.CloudManagerInterface {

	if wantedState.VMs != nil {
		l.econe.scaleVms(wantedState.VMs, l.kube)
		start := time.Now()
		isVmsDeployed := false
		for ; !isVmsDeployed && (time.Since(start) < deploymentTimeout); time.Sleep(10 * time.Second) {

			log.Info("waiting for VM to deploy")
			if areVMsCorrect(wantedState.VMs, l.econe.getVMs()) {
				log.Notice("Vms deployed")
				isVmsDeployed = true
			}

		}

		if !(time.Since(start) < deploymentTimeout) { //timeout exceed
			log.Warning("VM deployment timeout, moving on...")
		}

		start = time.Now()
		isKubernetesConfigured := false
		for ; !isKubernetesConfigured && (time.Since(start) < deploymentTimeout); time.Sleep(10 * time.Second) {

			nodesInKube := 0
			totalNumberofVMs := 0
			nodesList, err := l.kube.CoreV1().Nodes().List(metav1.ListOptions{}) //get node in kube
			if err != nil {
				panic(err)
			}
			for _, node := range nodesList.Items {
				if node.Status.Conditions[4].Status == apiv1.ConditionTrue { // 4 stands for isReady?
					nodesInKube++
				}
			}
			machines := l.econe.getVMs()
			for machineType := range machines {
				totalNumberofVMs += machines[machineType]
			}

			if nodesInKube != totalNumberofVMs {
				log.Debug("kube nodes:%v , vm number: %v\n", nodesInKube, totalNumberofVMs)
			} else {
				log.Info("Kubernetes configured, node count: %v", nodesInKube)
				isKubernetesConfigured = true
			}
			log.Info("waiting for VMs to join kubernetes")

		}
		if !(time.Since(start) < deploymentTimeout) { //timeout exceed
			log.Warning("Kubernetes join timeout, moving on...")
		}

	} else {
		log.Debug("%s has no VM state, keeping current configuration", wantedState.Name)
	}

	for key := range wantedState.Services {
		l.scaleContainers(key, wantedState.Services[key])
	}

	l.lastDeployedState = wantedState
	return l
}

//GetActiveState gets current application state
func (l Lrz) GetActiveState() ymlparser.State {
	return ymlparser.State{
		VMs:      l.econe.getVMs(),
		Services: l.getServiceCount(),
	}

}

//GetLastDeployedState returns last deployed state by the system
func (l Lrz) GetLastDeployedState() ymlparser.State {
	return l.lastDeployedState

}

//CheckState checks whether the deployed state and the actual state are the same
func (l Lrz) CheckState() bool {

	weDeployed := l.GetLastDeployedState()
	real := l.GetActiveState() //SORT

	//compare services

	//compare vms
	if areVMsCorrect(weDeployed.VMs, real.VMs) && areServicesCorrect(weDeployed.Services, real.Services) {

		return true
	}
	log.Error("ERROR:\ndepl: %#v\nreal: %#v", weDeployed, real)
	return false
}

func (l Lrz) getServiceCount() ymlparser.Service {

	deploymentList, _ := l.kube.AppsV1().Deployments(apiv1.NamespaceDefault).List(metav1.ListOptions{})

	currentServices := ymlparser.Service{}

	for _, d := range deploymentList.Items {
		currentServices[d.Name] = ymlparser.ServiceInfo{Replicas: int(*d.Spec.Replicas), CPU: d.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().Value(), Memory: d.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().Value()}
	}

	return currentServices
}

func (l Lrz) scaleContainers(serviceName string, serviceInfo ymlparser.ServiceInfo) string {

	log.Info("Updating Services...")
	deploymentsClient := l.kube.AppsV1().Deployments(apiv1.NamespaceDefault)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get(serviceName, metav1.GetOptions{})
		if getErr != nil {
			log.Critical("Failed to get latest version of Deployment: %v", getErr)

		}

		sn := int32(serviceInfo.Replicas)
		result.Spec.Replicas = &sn
		//result.Spec.Template.Spec.Containers[0].Args = []string{"-cpus", serviceInfo.CPU}
		result.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().Set(serviceInfo.Memory)
		result.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().Set(serviceInfo.CPU)

		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		log.Critical("Update failed: %v", retryErr)
	}
	log.Notice("Updated deployment...")

	return ""
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
