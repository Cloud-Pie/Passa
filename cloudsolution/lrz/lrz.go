package lrz

import (
	"fmt"
	"log"
	"sort"

	"github.com/Cloud-Pie/Passa/cloudsolution"
	"github.com/Cloud-Pie/Passa/ymlparser"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

type econe struct {
	username string
	password string
}
type Lrz struct {
	lastDeployedState   ymlparser.State
	econe               econe
	kube                *kubernetes.Clientset
	isActivelyDeploying bool
}

const ec2URL = "https://www.cloud.mwn.de:22"

//NewLRZManager return a new manager for lrz.
func NewLRZManager(username, password, configFile string) Lrz {

	config, _ := clientcmd.BuildConfigFromFlags("", configFile)
	clientset, _ := kubernetes.NewForConfig(config)
	cs := Lrz{
		econe: econe{
			username: username,
			password: password,
		},
		kube: clientset,
	}
	cs.lastDeployedState = cs.GetLastDeployedState()
	return cs
}

func (l Lrz) ChangeState(wantedState ymlparser.State) cloudsolution.CloudManagerInterface {

	l.isActivelyDeploying = true
	for _, service := range wantedState.Services {
		l.scaleContainers(service.Name, service.Scale)
	}
	l.lastDeployedState = l.GetActiveState()
	l.isActivelyDeploying = false
	return l
}

func (l Lrz) GetActiveState() ymlparser.State {
	return ymlparser.State{
		VMs:      l.getMachines(),
		Services: l.getServiceCount(),
	}

}
func (l Lrz) GetLastDeployedState() ymlparser.State {
	return ymlparser.State{}

}
func (l Lrz) CheckState() bool {
	return true
}

func (l Lrz) getServiceCount() []ymlparser.Service {

	deploymentList, _ := l.kube.AppsV1().Deployments(apiv1.NamespaceDefault).List(metav1.ListOptions{})

	currentServices := []ymlparser.Service{}

	for _, d := range deploymentList.Items {
		currentServices = append(currentServices, ymlparser.Service{
			Name:  d.Name,
			Scale: int(*d.Spec.Replicas),
		})
	}
	sort.Slice(currentServices, func(i, j int) bool {
		return currentServices[i].Name > currentServices[j].Name
	})

	return currentServices
}

func (l Lrz) getMachines() []ymlparser.VM {

	nodesList, err := l.kube.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	return []ymlparser.VM{{
		Type:  "myMachine",
		Scale: len(nodesList.Items),
	}}
}

func (l Lrz) scaleContainers(serviceName string, scaleNum int) string {

	log.Println("Updating deployment...")
	deploymentsClient := l.kube.AppsV1().Deployments(apiv1.NamespaceDefault)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get(serviceName, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		sn := int32(scaleNum)
		result.Spec.Replicas = &sn // reduce replica count

		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")

	return ""
}

func (e econe) createNewVM() {
	//Write the code for a new VM in LRZ in 'econe'
}
