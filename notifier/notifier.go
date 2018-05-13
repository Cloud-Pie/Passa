//Package notifier provides functions for Telegram communication.
package notifier

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/telegram-bot-api.v4"
)

var bot *tgbotapi.BotAPI
var chatID int64

//InitializeClient creates a new telegram bot
func InitializeClient() {
	telegramToken := os.Getenv("telegramToken")
	if telegramToken == "" {
		return
	}
	var err error
	chatID, err = strconv.ParseInt(os.Getenv("chatID"), 10, 64)
	if err != nil {
		return
	}
	bot, err = tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
}

//Notify notifies the user with the preffered channel
func Notify(message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	//bot.Send(msg)
	log.Println(msg)
}
