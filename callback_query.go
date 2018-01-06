package main

import (
	"fmt"
	"strings"

	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

const (
	menuLang   = "languages"
	menuRating = "ratings"
	menuRes    = "resources"

	switcherStatus = " 👈🏻"
)

var toggleStatus = map[bool]string{
	true:  "✅ ",
	false: "☑️ ",
}

func callbackQuery(call *tg.CallbackQuery) {
	usr, err := dbGetUserElseAdd(call.From.ID, call.From.LanguageCode)
	errCheck(err)

	data := strings.Split(call.Data, ":")
	switch data[0] {
	case "to":
		callbackTo(usr, call, data[1])
	case "switch":
		callbackSwitch(usr, call, data[1:])
	case "toggle":
		callbackToggle(usr, call, data[1:])
	case "add":
		callbackAdd(usr, call, data[1:])
	case "remove":
		callbackRemove(usr, call, data[1:])
	default:
		callbackAlert(call, "😱 Oh no!..")
	}
}

func callbackTo(usr *user, call *tg.CallbackQuery, option string) {
	log.Ln("callbackTo", option)
	switch option {
	case "settings":
		callbackToSettings(usr, call)
	case "languages":
		callbackToLanguages(usr, call)
	case "resources":
		callbackToResources(usr, call)
	case "ratings":
		callbackToRatings(usr, call)
	case "types":
		callbackToTypes(usr, call)
	case "blacklist":
		callbackToList(usr, call, blackList)
	case "whitelist":
		callbackToList(usr, call, whiteList)
	default:
		callbackAlert(call, "😱 Oh no!..")
	}
}

func callbackSwitch(usr *user, call *tg.CallbackQuery, options []string) {
	log.Ln("callbackSwitch", options[0])
	switch options[0] {
	case "language":
		callbackSwitchLanguage(usr, call, options[1])
	default:
		callbackAlert(call, "😱 Oh no!..")
	}
}

func callbackToggle(usr *user, call *tg.CallbackQuery, options []string) {
	log.Ln("callbackToggle", options[0])
	switch options[0] {
	case "resource":
		callbackToggleResource(usr, call, options[1])
	case "rating":
		callbackToggleRating(usr, call, options[1])
	case "type":
		callbackToggleTypes(usr, call, options[1])
	default:
		callbackAlert(call, "😱 Oh no!..")
	}
}

func callbackAdd(usr *user, call *tg.CallbackQuery, options []string) {
	switch options[0] {
	case "tags":
		go func() {
			_, err := bot.AnswerCallbackQuery(tg.NewAnswerCallbackQuery(call.ID))
			errCheck(err)
		}()

		usr, err := dbGetUserElseAdd(call.From.ID, call.From.LanguageCode)
		errCheck(err)

		T, err := langSwitch(usr.Language, call.From.LanguageCode)
		errCheck(err)

		text := T(fmt.Sprint("message_input_", options[1], "_tags"), map[string]interface{}{
			"Limit": 5,
		})

		reply := tg.NewMessage(int64(call.From.ID), text)
		reply.ParseMode = tg.ModeMarkdown
		reply.ReplyMarkup = tg.NewForceReply(true)

		_, err = bot.SendMessage(reply)
		errCheck(err)
	default:
		callbackAlert(call, "😱 Oh no!..")
	}
}

func callbackRemove(usr *user, call *tg.CallbackQuery, options []string) {
	switch options[0] {
	case blackList, whiteList:
		err := usr.removeListTag(options[0], strings.Join(options[1:], ""))
		errCheck(err)

		callbackUpdateListKeyboard(usr, call, options[0])
	default:
		callbackAlert(call, "😱 Oh no!..")
	}
}
