package main

import (
	"log"

	"map/handlers"
	"map/tgbotapi"
	t "map/types"
	"map/utils"
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

	// Loop through incoming updates
	for update := range updates {
		utils.Log(update)
		// Check if the update is a callback query
		if update.CallbackQuery != nil {
			// Check if the callback query is from a group chat
			switch update.CallbackQuery.Message.Chat.Type {
			case "group":
				// If the callback data is "connect", connect the user to the chat map
				if update.CallbackQuery.Data == "connect" {
					handlers.ConnectUserToChatMap(update, bot)
					continue // Continue to the next update
				}
			case "private":
				// If the callback query is from a private chat, continue to the next update
				continue
			}
		}

		// Check if the update is a message
		if update.Message != nil {
			// Log the message

			// If the message is from a group chat and a new member joined
			if (update.Message.NewChatMembers != nil) && (update.Message.Chat.Type == "group") {
				// If the new member is the bot, connect it to the chat
				if (update.Message.NewChatMembers[0].ID == bot.Self.ID) && (update.Message.Chat.Type == "group") {
					handlers.ConnectBotToChat(update, bot)
					continue // Continue to the next update
				}
			}

			// Check the type of chat the message is from
			switch update.Message.Chat.Type {
			case "private":
				// If the message is from a private chat and the text is "/start", start the bot
				if update.Message.Text == "/start" {
					handlers.StartBot(update, bot)
					continue // Continue to the next update
				}
				// If the message is from a private chat, create a spot
				handlers.CreateSpot(update, bot, activeSpotCreations)
			case "group":
				// If the message is from a group chat, handle the message in the group
				handlers.MessageInGroup(update, bot)
			}
		}
	}
}
