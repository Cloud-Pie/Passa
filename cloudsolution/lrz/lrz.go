package lrz

import (
	"fmt"
	"log"
	"reflect"
	"sort"

	"github.com/Cloud-Pie/Passa/cloudsolution"
	"github.com/Cloud-Pie/Passa/ymlparser"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

//Lrz keeps the data needed for econe and kubernetes interfaces.
type Lrz struct {
	lastDeployedState   ymlparser.State
	econe               econe
	kube                *kubernetes.Clientset
	isActivelyDeploying bool
}

//NewLRZManager return a new manager for lrz.
func NewLRZManager(username, password, configFile string) Lrz {

	config, err := clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		log.Fatal("Cannot not connect to kubernetes cluster , exiting...")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Cannot not connect to kubernetes cluster , exiting...")
	}
	cs := Lrz{
		econe: econe{
			username: username,
			password: password,
		},
		kube: clientset,
	}
	cs.lastDeployedState = cs.GetActiveState()
	return cs
}

//ChangeState deploys wanted state to LRZ and kubernetes
func (l Lrz) ChangeState(wantedState ymlparser.State) cloudsolution.CloudManagerInterface {

	l.isActivelyDeploying = true

	if wantedState.VMs != nil {
		l.econe.scaleVms(wantedState.VMs, l.GetLastDeployedState().VMs)
	} else {
		log.Printf("%s has no VM state, keeping current...", wantedState.Name)
	}
	for _, service := range wantedState.Services {
		l.scaleContainers(service.Name, service.Scale)
	}
	l.lastDeployedState = l.GetActiveState()
	l.isActivelyDeploying = false
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
	if l.isActivelyDeploying { //BUG: This doesn't read with mutex, we will give wrong error eventually.
		log.Println("Actively deploying new state")
		return true
	}
	weDeployed := l.GetLastDeployedState()
	real := l.GetActiveState() //SORT

	sort.Slice(weDeployed.Services, func(i, j int) bool {
		return weDeployed.Services[i].Name > weDeployed.Services[j].Name
	})

	sort.Slice(real.Services, func(i, j int) bool {
		return real.Services[i].Name > real.Services[j].Name
	})

	real.ISODate = weDeployed.ISODate //see dockerswarm.go
	if reflect.DeepEqual(weDeployed, real) {
		log.Println("State holds")
		return true
	}

	log.Printf("ERROR: deployed: %#v real: %#v", weDeployed, real)
	return false

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
