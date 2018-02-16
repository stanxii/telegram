package i18n

import "github.com/nicksnyder/go-i18n/i18n"

const Fallback = "en"

func SwitchTo(codes ...string) (T i18n.TranslateFunc, err error) {
	codes = append(codes, Fallback)
	T, err = i18n.Tfunc(codes[0], codes[1:]...)
	return
}