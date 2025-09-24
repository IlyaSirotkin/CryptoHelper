package telegram_service

import (
	"cryptoHelper/internal/datasource/datasource_interface"
	"cryptoHelper/internal/display/display_interface"
	"cryptoHelper/internal/display/telegram_display"
	logger "cryptoHelper/pkg/applogger"
	"errors"
	"fmt"
	"os"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	botAPI     *tgBotAPI.BotAPI
	datasource datasource_interface.Datasource
	display    display_interface.Display
}

func NewTelegram(token string) (*Telegram, error) {
	bot, err := tgBotAPI.NewBotAPI(os.Getenv(token))
	if err != nil {
		logger.Get().Error("Telegram wasn't created NewBotAPI return err")
		return nil, err
	} else {
		logger.Get().Debug("Telegram successfully created botAPI")
		return &Telegram{botAPI: bot}, nil
	}
}

func (t Telegram) SetInput(dsrc datasource_interface.Datasource) error {
	if dsrc == nil {
		logger.Get().Error("Datasource is nil in Telegram SetInput")
		return errors.New("datasource is nil in Telegram SetInput")
	} else {
		t.datasource = dsrc
		logger.Get().Debug("Telegram service was successfully set input device")
		return nil
	}
}

func (t Telegram) SetOutput(dspl display_interface.Display) error {
	if dspl == nil {
		logger.Get().Error("Display is nil in Telegram SetOutput")
		return errors.New("display is nil in Telegram SetOutput")
	} else {
		t.display = dspl
		logger.Get().Debug("Telegram service was successfully set output device")
		return nil
	}
}

func (t Telegram) GetData(currencyName string) (float32, error) {
	if t.datasource != nil {
		logger.Get().Debug("Datasource_handler called ExtractCurrentPrice() successfully")
		return t.datasource.ExtractCurrentPrice(currencyName)
	} else {
		logger.Get().Error("Datasource_interface is nil, GetData() operation cannot be completed")
		return 0.0, errors.New("datasource_interface is nil, operation can not be completed")
	}
}

func (t Telegram) SendData(message string) error {
	if t.display != nil {
		err := t.display.SendMessage(message)
		if err != nil {
			return fmt.Errorf("SendMessage return error%w", err)
		}
		logger.Get().Debug("Display_handler called SendData() successfully")
		return nil
	} else {
		logger.Get().Error("Display_interface is nil, SendData() operation cannot be completed")
		return errors.New("display_interface is nil, operation can not be completed")
	}
}

func (t *Telegram) Update() {
	updateConfig := tgBotAPI.NewUpdate(0)
	updateConfig.Timeout = 60

	updateChan := t.botAPI.GetUpdatesChan(updateConfig)

	for update := range updateChan {
		if update.Message != nil {
			chatID := update.Message.Chat.ID

			telegram_display.SetChatID(chatID, t.display.(*telegram_display.Bot))
			text := update.Message.Text

			switch text {
			case "/start":

			case "/help":

			default:

			}
		}
	}
}
