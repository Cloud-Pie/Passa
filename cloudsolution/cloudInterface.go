package cloudsolution

import "github.com/Cloud-Pie/Passa/ymlparser"

//CloudManagerInterface is the interface for cloud management system
type CloudManagerInterface interface {
	ChangeState(ymlparser.State) []string
}
