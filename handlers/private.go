package handlers

import (
	"context"
	"log"
	"map/db"
	"map/helpers"
	"map/tgbotapi"
	t "map/types"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
)

func StartBot(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	clientOptions := db.GetClientOptions()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Creating your account..."
	bot.Send(msg)

	// create a new user
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	client.Database("data").Collection("users").InsertOne(context.TODO(), t.User{
		Username: update.Message.From.UserName,
		ID:       update.Message.From.ID,
		Chats:    []int64{},
	})
	

	// send a success message
	msg.Text = "Account created!"
	msg.ReplyMarkup = helpers.GetWebAppKeyboard("1=2")
	bot.Send(msg)
}

func CreateSpot(update tgbotapi.Update, bot *tgbotapi.BotAPI, activeSpotCreations t.ActiveSpotCreations) {

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyMarkup = helpers.PrivateChatMainKeyboard
	msg.ParseMode = "Markdown"

	activeCreationSession, ok := activeSpotCreations[update.Message.Chat.ID]

	if ok {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		switch activeSpotCreations[update.Message.Chat.ID].Step {
		case 0:
			activeCreationSession.Data.Title = update.Message.Text
			activeCreationSession.Step++
			activeSpotCreations[update.Message.Chat.ID] = activeCreationSession
			msg.Text = "Please enter a description"
		case 1:
			activeCreationSession.Data.Description = update.Message.Text
			activeCreationSession.Step++
			activeSpotCreations[update.Message.Chat.ID] = activeCreationSession
			msg.Text = "Please enter a category"
		case 2:
			activeCreationSession.Data.Category = update.Message.Text
			activeCreationSession.Step++
			activeSpotCreations[update.Message.Chat.ID] = activeCreationSession
			msg.Text = "Please enter a location"
		case 3:
			if update.Message.Location != nil {
				clientOptions := db.GetClientOptions()

				activeCreationSession.Data.Lat = update.Message.Location.Latitude
				activeCreationSession.Data.Lon = update.Message.Location.Longitude
				activeCreationSession.Step++
				activeSpotCreations[update.Message.Chat.ID] = activeCreationSession
				// msg.Text = "Please enter a picture"

				delete(activeSpotCreations, update.Message.Chat.ID)

				msg.Text = "Creating spot..."
				bot.Send(msg)

				// connect to mongodb
				client, err := mongo.Connect(context.TODO(), clientOptions)
				if err != nil {
					log.Fatal(err)
				}

				// get collection and create new message
				coll := client.Database("points").Collection(strconv.FormatInt(update.Message.Chat.ID, 10))

				newSpot := activeCreationSession.Data
				result, err := coll.InsertOne(context.TODO(), newSpot)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("Create a spot: ", result.InsertedID)

				msg.ReplyMarkup = helpers.PrivateChatMainKeyboard
				msg.Text = "Spot created!\n\n*Title:* " + newSpot.Title + "\n*Description:* " + newSpot.Description + "\n*Category:* " + newSpot.Category
			} else {
				msg.Text = "Please send a Telegram location message"
			}
		case 4:
			// activeCreationSession.Data.Pic = update.Message.Text
			// activeSpotCreations[update.Message.Chat.ID] = activeCreationSession
			// remove active spot creation

			// delete(activeSpotCreations, update.Message.Chat.ID)

			// msg.Text = "Creating spot..."
			// bot.Send(msg)

			// // connect to mongodb
			// ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			// defer cancel()
			// client, err := mongo.Connect(ctx, clientOptions)
			// if err != nil {
			// 	log.Fatal(err)
			// }

			// // get collection and create new message
			// coll := client.Database("points").Collection(strconv.FormatInt(update.Message.Chat.ID, 10))

			// newSpot := activeCreationSession.Data
			// result, err := coll.InsertOne(ctx, newSpot)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// log.Println("Create a spot: ", result.InsertedID)

			// msg.ReplyMarkup = helpers.PrivateChatMainKeyboard
			// msg.Text = "Spot created!\n\n*Title:* " + newSpot.Title + "\n*Description:* " + newSpot.Description + "\n*Category:* " + newSpot.Category
		}

	} else {
		switch update.Message.Text {
		case "add spot":
			{
				activeSpotCreations[update.Message.Chat.ID] = t.ActiveSpotCreation{
					ChatID: update.Message.Chat.ID,
					Step:   0,
					Data: t.Spot{
						Author: update.Message.From.UserName,
					},
				}
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				msg.Text = "Please enter a title"
			}
		case "get all chats":
			{
				msg.Text = "Getting chats..."
			}
		default:
			{
				msg.Text = "Choose an action"
			}
		}
	}
	bot.Send(msg)
}
