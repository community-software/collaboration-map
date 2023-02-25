package main

import (
	"encoding/json"
	"log"

	"map/handlers"
	"map/tgbotapi"
	t "map/types"
)

func main() {
	// create a empty map of active spot creations
	activeSpotCreations := make(t.ActiveSpotCreations)

	// create a bot
	bot, err := tgbotapi.NewBotAPI("6120325657:AAEFiqq9g84YR1-tzfPbf-zDvcbgk2zqQ90")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// listen for updates
	for update := range updates {
		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Message.Chat.Type {
			case "group":
				{
					if update.CallbackQuery.Data == "connect" {
						handlers.ConnectUserToChatMap(update, bot)
						continue
					}
				}
			case "private":
				continue
			}
		}
		if update.Message != nil { // If we got a message
			Log(update)

			if (update.Message.NewChatMembers != nil) && (update.Message.Chat.Type == "group") {
				if (update.Message.NewChatMembers[0].ID == bot.Self.ID) && (update.Message.Chat.Type == "group") {
					handlers.ConnectBotToChat(update, bot)
					continue
				}
			}

			if (update.Message.Chat.Type == "private") && (update.Message.Text == "/start") {
				handlers.StartBot(update, bot)
				continue
			}

			switch update.Message.Chat.Type {
			case "private":
				handlers.CreateSpot(update, bot, activeSpotCreations)
			case "group":
				handlers.MessageInGroup(update, bot)
			}
		}
	}
}

func Log(u tgbotapi.Update) {
	newMsg, _ := json.Marshal(u)

	log.Printf("\033[1;34m%s\033[0m ", "Message received: ")
	log.Printf("\033[1;34m%s\033[0m ", string(newMsg))
	log.Printf("\033[1;34m%s\033[0m ", "End message")
}
