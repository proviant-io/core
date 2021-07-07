package i18n

import (
	"fmt"
	"log"
)

type Locale string

const En = "en"
const Ru = "ru"
const Es = "es"

type Message struct {
	Template string
	Params   []interface{}
}

func NewMessage(template string, params ...interface{}) Message {
	return Message{Template: template, Params: params}
}

type Localizer interface {
	T(Message, Locale) string
	Missing() []string
}

type FileLocalizer struct {
	strings map[string]map[Locale]string
	missing []string
}

func (l *FileLocalizer) Missing() []string {
	return l.missing
}

func (l *FileLocalizer) T(m Message, locale Locale) string {

	if translations, ok := l.strings[m.Template]; ok {
		if translation, ok := translations[locale]; ok {
			if len(m.Params) > 0 {
				log.Println(translation, m.Params)
				return fmt.Sprintf(translation, m.Params...)
			}
			return translation

		}
	}

	log.Printf("missing translation: %s - %s \n", locale, m.Template)
	l.missing = append(l.missing, m.Template)

	if len(m.Params) > 0 {
		log.Println(m.Template, m.Params)
		return fmt.Sprintf(m.Template, m.Params...)
	}
	return m.Template
}

func NewFileLocalizer() Localizer {

	strings := map[string]map[Locale]string{
		"product with id %d not found": {
			En: "product with id %d not found",
			Ru: "продукт под номером %d не найден",
		},
		"category with id %d not found": {
			En: "category with id %d not found",
			Ru: "категория под номером %d не найдена",
		},
		"list with id %d not found": {
			En: "list with id %d not found",
			Ru: "список под номером %d не найден",
		},
		"stock with id %d not found": {
			En: "stock with id %d not found",
			Ru: "запись остатка под номером %d не найдена",
		},
		"id cannot be empty": {
			En: "id cannot be empty",
			Ru: "ID не может быть пустым",
		},
		"id is not a number: %v": {
			En: "id is not a number: %v",
			Ru: "ID не является числом: %v",
		},
		"parse payload error: %v": {
			En: "parse payload error: %v",
			Ru: "ошибка обработки данных запроса: %v",
		},
		"title should not be empty": {
			En: "title should not be empty",
			Ru: "заголовок должен быть заполнен",
		},
		"You can't remove list with products. Clean products first.": {
			En: "You can't remove list with products. Clean products first.",
			Ru: "Вы не можете удалить список с продуктами. Удалите или перенесите продукты перед удалением.",
		},
	}

	return &FileLocalizer{strings, []string{}}
}

func LocaleFromString(s string) Locale {
	switch s {
	case "en":
		return En
	case "ru":
		return Ru
	default:
		// TODO log that locale is not supported
		return En
	}
}
