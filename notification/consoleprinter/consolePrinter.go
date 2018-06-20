//Package consoleprinter provides functions for Telegram communication.
package consoleprinter

import (
	"github.com/op/go-logging"
)

type consolePrinter struct{}

var log = logging.MustGetLogger("passa")

//InitializeClient starts the console client
func InitializeClient() consolePrinter {
	return consolePrinter{}
}

//Notify notifies the user with the preffered channel
func (cp consolePrinter) Notify(message string) {
	log.Info(message)
}
