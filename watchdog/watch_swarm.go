package watchdog

import "os/exec"

func getStatus() {
	exec.Command("docker service ls --format '{{.Name}} {{.Replicas}}'")
}
