package cloudsolution

import (
	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
)

//CloudManagerInterface is the interface for cloud management system
type CloudManagerInterface interface {
	ChangeState(ymlparser.State) []string
}
