package telegram_display

import (
	logger "cryptoHelper/pkg/applogger"
	"fmt"
	"os"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SetSenderChatID(id int64, b *BotSender) {
	b.chatID = id
}

type BotSender struct {
	botAPI *tgBotAPI.BotAPI
	chatID int64
}

func NewBotSender(token string) (*BotSender, error) {
	bot, err := tgBotAPI.NewBotAPI(os.Getenv(token))
	if bot != nil && err == nil {
		logger.Get().Debug("NewBotSender successfully created BotSender")
		return &BotSender{botAPI: bot}, nil
	} else {
		logger.Get().Error("tgBotAPI NewBotAPI() return error " + fmt.Sprint(err))
		return nil, fmt.Errorf("tgBotAPI is nil in Bot constructor%w", err)
	}
}

func (t BotSender) SendMessage(message string) error {
	msg := tgBotAPI.NewMessage(t.chatID, message)
	msg.ParseMode = "Markdown"
	_, err := t.botAPI.Send(msg)
	if err != nil {
		logger.Get().Error("Send() return error in BotSender::SendMessage func " + fmt.Sprint(err))
		return fmt.Errorf("send() return error in SendMessage func %w", err)
	} else {
		logger.Get().Debug("Send() in BotSender successfully send message")
		return nil
	}
}

type BotMarkupSender struct {
	botAPI *tgBotAPI.BotAPI
	chatID int64
}

func SetMarkupSenderChatID(id int64, b *BotMarkupSender) {
	b.chatID = id
}

func NewBotMarkupSender(token string) (*BotMarkupSender, error) {
	bot, err := tgBotAPI.NewBotAPI(os.Getenv(token))
	if bot != nil && err == nil {
		logger.Get().Debug("NewBotMarkupSender successfully created BotMarkupSender")
		return &BotMarkupSender{botAPI: bot}, nil
	} else {
		logger.Get().Error("tgBotAPI NewBotAPI() return error " + fmt.Sprint(err))
		return nil, fmt.Errorf("tgBotAPI is nil in Bot constructor%w", err)
	}
}

func (t BotMarkupSender) SendMessage(message string) error {
	msg := tgBotAPI.NewMessage(t.chatID, message)

	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData("BTC", "btc_section"),
			tgBotAPI.NewInlineKeyboardButtonData("ETH", "eth_section"),
			tgBotAPI.NewInlineKeyboardButtonData("SOL", "sol_section"),
			tgBotAPI.NewInlineKeyboardButtonData("ADA", "ada_section"),
			tgBotAPI.NewInlineKeyboardButtonData("XRP", "xrp_section"),
			tgBotAPI.NewInlineKeyboardButtonData("ONDO", "ondo_section"),
		),
	)
	_, err := t.botAPI.Send(msg)
	if err != nil {
		logger.Get().Error("Send() return error in BotMarkupSender::SendMessage func " + fmt.Sprint(err))
		return fmt.Errorf("Send() return error in BotMarkupsender::SendMessage func %w", err)
	} else {
		logger.Get().Debug("Send() in BotMarkupSender successfully send message")
		return nil
	}
}
