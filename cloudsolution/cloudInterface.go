package cloudsolution

import "github.com/Cloud-Pie/Passa/ymlparser"

//CloudManagerInterface is the interface for cloud management system
type CloudManagerInterface interface {
	ChangeState(ymlparser.State) CloudManagerInterface
	GetActiveState() ymlparser.State
	GetLastDeployedState() ymlparser.State
	CheckState() bool
}
