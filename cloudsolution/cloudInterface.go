package cloudsolution

import (
	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
)

//CloudManager is the interface for cloud management system
type CloudManager interface {
	ChangeState(ymlparser.Service) []string
}
