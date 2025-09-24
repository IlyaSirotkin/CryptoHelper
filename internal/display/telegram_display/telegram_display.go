package telegram_display

import (
	"fmt"
	"os"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SetChatID(id int64, b *Bot) {
	b.chatID = id
}

type Bot struct {
	botAPI *tgBotAPI.BotAPI
	chatID int64
}

func NewBot(token string) (*Bot, error) {
	bot, err := tgBotAPI.NewBotAPI(os.Getenv(token))
	if bot != nil && err == nil {
		return &Bot{botAPI: bot}, nil
	} else {
		return nil, fmt.Errorf("tgBotAPI is nil in Bot constructor%w", err)
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
