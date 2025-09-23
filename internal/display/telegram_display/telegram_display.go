package telegram_display

import (
	"errors"
	"fmt"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	chatID int64
	botAPI *tgBotAPI.BotAPI
}

func NewTelegram(bot *tgBotAPI.BotAPI, id int64) (*Bot, error) {
	if bot != nil {
		return &Bot{botAPI: bot, chatID: id}, nil
	} else {
		return nil, errors.New("tgBotAPI is nil in Bot constructor")
	}

}

func (t Bot) SendMessage(message string) error {
	msg := tgBotAPI.NewMessage(t.chatID, message)
	msg.ParseMode = "Markdown"
	_, err := t.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("Send() return error in SendMessage func %w", err)
	} else {
		return nil
	}
}
