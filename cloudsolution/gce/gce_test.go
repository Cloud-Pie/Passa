package gce

import (
	"testing"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

func Test_a(t *testing.T) {
	g := NewGCEManager("hpa-cluster")
	ss := ymlparser.Service{"movieapp": ymlparser.ServiceInfo{
		Replicas: 3,
		CPU:      "700m",
		Memory:   700000000,
	}}
	g.scaleContainers(ss)
}
