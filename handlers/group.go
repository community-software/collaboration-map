package handlers

import (
	"context"
	"log"
	"map/db"
	"map/helpers"
	"map/tgbotapi"
	t "map/types"
	"map/utils"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConnectUserToChatMap(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	clientOptions := db.GetClientOptions()
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := bot.Request(callback); err != nil {
		log.Fatal(err)
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	var userObj t.User
	findErr := client.Database("data").Collection("users").FindOne(context.TODO(), bson.M{"id": update.CallbackQuery.From.ID}).Decode(&userObj)
	if findErr != nil {
		client.Database("data").Collection("users").InsertOne(context.TODO(), t.User{
			Username: update.CallbackQuery.From.UserName,
			ID:       update.CallbackQuery.From.ID,
			Chats:    []int64{update.CallbackQuery.Message.Chat.ID},
		})
		msg.Text = "@" + update.CallbackQuery.From.UserName + ", you are on map!\nPlease, add me to your private chat and push start"
	} else {
		chats := userObj.Chats

		if utils.Contains(chats, update.CallbackQuery.Message.Chat.ID) {
			msg.Text = "@" + update.CallbackQuery.From.UserName + ", you are already on map!"
		} else {
			client.Database("data").Collection("users").UpdateOne(
				context.TODO(),
				bson.M{"id": update.CallbackQuery.From.ID},
				bson.D{{
					Key: "$set",
					Value: bson.D{{
						Key: "chats", Value: append(userObj.Chats, update.CallbackQuery.Message.Chat.ID),
					}},
				}},
			)
			msg.Text = "@" + update.CallbackQuery.From.UserName + ", you are on map!\nPlease, add me to your private chat and push start"
		}

	}
	bot.Send(msg)
}

func ConnectBotToChat(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	clientOptions := db.GetClientOptions()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Creating a chat map..."
	bot.Send(msg)

	// create a new chat map
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	client.Database("points").CreateCollection(context.TODO(), strconv.FormatInt(update.Message.Chat.ID, 10))
	client.Database("data").Collection("chats").InsertOne(context.TODO(), t.Chat{
		ID:    update.Message.Chat.ID,
		Title: update.Message.Chat.Title,
	})

	msg.Text = "Chat map created!\n\nIf you want to add a spots and use the map, you need to click on connect button, and add me to your private chat"
	msg.ReplyMarkup = helpers.GroupChatMainKeyboard
	bot.Send(msg)
}

func MessageInGroup(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	msg.ReplyMarkup = helpers.GroupChatMainKeyboard

	msg.Text = "Its a group chat"
	bot.Send(msg)
}
