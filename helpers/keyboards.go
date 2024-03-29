package helpers

import "map/tgbotapi"

func GetWebAppKeyboard(params string) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("add spot"),
			tgbotapi.NewKeyboardButtonWebApp("show spots", tgbotapi.WebAppInfo{
				URL: "https://collaboration-app-webapp.vercel.app/?" + params,
			}),
			tgbotapi.NewKeyboardButton("get all chats"),
		),
	)
}

var PrivateChatMainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("add spot"),
		tgbotapi.NewKeyboardButtonWebApp("show spots", tgbotapi.WebAppInfo{
			URL: "https://collaboration-app-webapp.vercel.app/",
		}),
		tgbotapi.NewKeyboardButton("get all chats"),
	),
)
var GroupChatMainKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("connect", "connect"),
	),
)
