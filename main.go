package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var spamMessages []string

func handleErrorAndExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func checkMessageForSpam(spamMessages []string, messageText string) bool {
	for _, spam := range spamMessages {
		if strings.TrimSpace(spam) == strings.TrimSpace(messageText) {
			return true
		}
	}
	return false
}

func main() {
	// Read spam db and load it in RAM
	spamMessagesCSV, err := os.Open("spamMessages.csv")
	handleErrorAndExit(err)

	records := csv.NewReader(spamMessagesCSV)
	for {
		spam, err := records.Read()
		if err == io.EOF {
			break
		}
		spamMessages = append(spamMessages, spam[0])
	}

	// Create connection to telegram api
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	handleErrorAndExit(err)

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updateChannel := bot.GetUpdatesChan(updateConfig)

	// Each new message is checked. Removed if it matches the spam database
	for update := range updateChannel {
		if update.Message == nil {
			continue
		}
		shouldBeDeleted := checkMessageForSpam(spamMessages, update.Message.Text)
		if shouldBeDeleted {
			messageToDelete := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
			_, err := bot.Send(messageToDelete)
			handleError(err)
		}
	}

}
