//Package consoleprinter provides functions for Telegram communication.
package consoleprinter

import (
	"log"
)

type consolePrinter struct{}

func InitializeClient() consolePrinter {
	return consolePrinter{}
}

//Notify notifies the user with the preffered channel
func (cp consolePrinter) Notify(message string) {
	log.Println(message)
}
