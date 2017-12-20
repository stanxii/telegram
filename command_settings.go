package main

import (
	"strings"

	tg "github.com/toby3d/go-telegram"
)

func commandSettings(msg *tg.Message) {
	usr, err := dbGetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
	errCheck(err)

	_, err = bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	T, err := langSwitch(usr.Language, msg.From.LanguageCode)
	errCheck(err)

	var activeRes []string
	for k, v := range usr.Resources {
		if v {
			title := resources[k]["title"].(string)
			activeRes = append(activeRes, title)
		}
	}

	ratings, err := usr.getRatingsStatus()
	errCheck(err)

	text := T("message_settings", map[string]interface{}{
		"Language":  languageNames[usr.Language],
		"Resources": strings.Join(activeRes, ", "),
		"Ratings":   ratings,
		"Blacklist": strings.Join(usr.Blacklist, "`, `"),
		"Whitelist": strings.Join(usr.Whitelist, "`, `"),
	})

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T(
					"button_language",
					map[string]interface{}{"Flag": T("language_flag")},
				),
				"to:languages",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_resources"), "to:resources",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_ratings"), "to:ratings",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_blacklist"), "to:blacklist",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_whitelist"), "to:whitelist",
			),
		),
	)

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
