//Package database provides functions for database
package database

import "github.com/Cloud-Pie/Passa/ymlparser"

//SearchQuery returns the index of the state in config file
func SearchQuery(currentStates []ymlparser.State, searchName string) int {

	for idx := range currentStates {
		if currentStates[idx].Name == searchName {
			return idx
		}
	}
	return -1
}
