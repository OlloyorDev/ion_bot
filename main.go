package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6904390973:AAGMnOhTDYvJXOMFzsRnoPYeQP63jPiPKBM")
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err) // Improved error handling
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			// Log the received message to the console
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// Echo the message back to the user
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			log.Printf("INPUT : ", msg)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}
