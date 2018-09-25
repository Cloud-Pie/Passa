package gce

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Cloud-Pie/Passa/cloudsolution"
	"github.com/Cloud-Pie/Passa/ymlparser"
	"github.com/op/go-logging"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
)

//gcloud container clusters resize [CLUSTER_NAME] --node-pool [POOL_NAME] --size [SIZE]
const resizeClusterCommand = "gcloud container clusters  resize %s --node-pool %s --size %d -q"

var types = []string{"t2.micro", "t2.large"}

const getAccount = "gcloud info --format='value(config.account)'"

var log = logging.MustGetLogger("passa")

//GCE keeps configuration for Google Cloud Engine
type GCE struct {
	lastDeployedState ymlparser.State
	clusterName       string
	kube              *kubernetes.Clientset
}

//NewGCEManager return a new manager for GCE.
func NewGCEManager(cn string) GCE {
	if isCommandAvailable("gcloud") {
		log.Info("Commands found: gcloud")
	} else {

		log.Critical("gcloud not found")

	}
	accountName, _ := exec.Command("sh", "-c", getAccount).Output()
	log.Debug("Authenticated as: %s", string(accountName))

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		log.Fatal("Couldn't build config from file")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Cannot connect to kubernetes cluster , exiting...")
	}
	cs := GCE{
		clusterName: cn,
		kube:        clientset,
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

func (g GCE) scaleContainers(wantedContainers ymlparser.Service) string {

	for serviceName := range wantedContainers {
		log.Info("Updating Services...")
		deploymentsClient := g.kube.AppsV1().Deployments(apiv1.NamespaceDefault)
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Deployment before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			result, getErr := deploymentsClient.Get(serviceName, metav1.GetOptions{})
			if getErr != nil {
				log.Critical("Failed to get latest version of Deployment: %v", getErr)

			}

			sn := int32(wantedContainers[serviceName].Replicas)
			result.Spec.Replicas = &sn
			//result.Spec.Template.Spec.Containers[0].Args = []string{"-cpus", serviceInfo.CPU}

			result.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().Set(wantedContainers[serviceName].Memory)
			cpuInt64, _ := strconv.ParseInt(wantedContainers[serviceName].CPU, 10, 64)
			result.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().Set(cpuInt64)

			//kubectl patch deployment movieapp  --type json -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/resources/limits/memory", "value":"12312321313213"}]' --type json -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/resources/limits/cpu", "value":"200m"}]'

			//fmt.Println(cp.String())
			_, updateErr := deploymentsClient.Update(result)

			return updateErr
		})
		if retryErr != nil {
			log.Critical("Update failed: %v", retryErr)
		}
		log.Notice("Updated deployment...")
	}
	return ""
}

func (g GCE) getVMs() ymlparser.VM {
	vms := ymlparser.VM{}

	nodesList, err := g.kube.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Fatal("Couldn't build config from file")
	}
	for _, node := range nodesList.Items {
		for _, t := range types {
			searchString := strings.Replace(t, ".", "-", -1)

			if strings.Contains(node.Name, searchString) {
				vms[t]++
			}
		}
	}
	fmt.Println(vms)
	return vms

}

func (g GCE) getServices() ymlparser.Service {
	deploymentList, _ := g.kube.AppsV1().Deployments(apiv1.NamespaceDefault).List(metav1.ListOptions{})

	serviceMap := ymlparser.Service{}

	for _, d := range deploymentList.Items {

		serviceMap[d.Name] = ymlparser.ServiceInfo{Replicas: int(*d.Spec.Replicas),
			CPU:    d.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String(),
			Memory: d.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().MilliValue()}
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
