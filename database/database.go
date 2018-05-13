//Package database provides functions for database
package database

import "gitlab.lrz.de/ga53lis/PASSA/ymlparser"

//SearchQuery returns the index of the state in config file
func SearchQuery(currentStates []ymlparser.State, searchName string) int {

	for idx := range currentStates {
		if currentStates[idx].Name == searchName {
			return idx
		}
	}
	return -1
}
