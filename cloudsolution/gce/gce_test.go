package gce

import (
	"testing"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

func Test_a(t *testing.T) {
	g := NewGCEManager("hpa-cluster")
	ss := ymlparser.Service{"movieapp": ymlparser.ServiceInfo{
		Replicas: 5,
		CPU:      "200m",
		Memory:   100000000000,
	},
	}

	g.scaleContainers(ss)
	//kubectl patch deployment movieapp  --type json -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/resources/limits/memory", "value":"12312321313213"}]' --type json -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/resources/limits/cpu", "value":"200m"}]'
}
