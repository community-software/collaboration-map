package handlers

import (
	"context"
	"encoding/json"
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

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := bot.Request(callback); err != nil {
		log.Fatal(err)
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// msgStr, _ := json.Marshal(update.CallbackQuery.Message)
	// log.Println(msgStr)

	user, err := client.Database("data").Collection("users").FindOne(context.TODO(), bson.M{"id": update.CallbackQuery.From.ID}).DecodeBytes()
	if err != nil {
		log.Println("User was created")
		log.Println(err)
		client.Database("data").Collection("users").InsertOne(context.TODO(), t.User{
			Username: update.CallbackQuery.From.UserName,
			ID:       update.CallbackQuery.From.ID,
			Chats:    []int64{update.CallbackQuery.Message.Chat.ID},
		})
	} else {
		log.Println("User was found")
		var userObj t.User
		json.Unmarshal(user, &userObj)
		utils.Log(userObj)
		userObj.Chats = append(userObj.Chats, update.CallbackQuery.Message.Chat.ID)
		client.Database("data").Collection("users").UpdateOne(context.TODO(), t.User{ID: update.CallbackQuery.From.ID}, userObj)
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
	msg.Text = "@" + update.CallbackQuery.From.UserName + ", you are on map!\nPlease, add me to your private chat and push start"
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
