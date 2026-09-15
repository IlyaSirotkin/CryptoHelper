package telegram_service

import (
	"context"
	cacheinterface "cryptoHelper/internal/cache/cache_interface"
	"cryptoHelper/internal/datasource/datasource_interface"
	"cryptoHelper/internal/display/display_interface"
	"cryptoHelper/internal/display/telegram_display"
	logger "cryptoHelper/pkg/applogger"
	"errors"
	"fmt"
	"os"
	"strconv"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

const NotExist = redis.Nil

type Telegram struct {
	botAPI      *tgBotAPI.BotAPI
	datasource  datasource_interface.Datasource
	display     display_interface.Display
	swapDisplay display_interface.Display
	cache       cacheinterface.CacheHandler
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

func (t *Telegram) SetInput(dsrc datasource_interface.Datasource) error {
	if dsrc == nil {
		logger.Get().Error("Datasource is nil in Telegram SetInput")
		return errors.New("datasource is nil in Telegram SetInput")
	} else {
		t.datasource = dsrc
		logger.Get().Debug("Telegram service was successfully set input device")
		return nil
	}
}

func (t *Telegram) SetCache(cache cacheinterface.CacheHandler) error {
	if cache == nil {
		logger.Get().Error("Cache is nil in Telegram SetCache")
		return errors.New("Cache is nil in Telegram SetCache")
	} else {
		t.cache = cache
		logger.Get().Debug("Telegram service was successfully set cache device")
		return nil
	}
}

func (t *Telegram) SetOutput(dspl display_interface.Display) error {
	if dspl == nil {
		logger.Get().Error("Display is nil in Telegram SetOutput")
		return errors.New("display is nil in Telegram SetOutput")
	} else {
		t.display = dspl
		logger.Get().Debug("Telegram service was successfully set output device")
		return nil
	}
}

func (t Telegram) GetData(currencyName string) (float64, error) {
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
			return fmt.Errorf("SendMessage return error %w", err)
		}
		logger.Get().Debug("Display_handler called SendData() successfully")
		return nil
	} else {
		logger.Get().Error("Display_interface is nil, SendData() operation cannot be completed")
		return errors.New("display_interface is nil, operation can not be completed")
	}
}

func (t *Telegram) Update() error {

	updateConfig := tgBotAPI.NewUpdate(0)
	updateConfig.Timeout = 60

	updateChan := t.botAPI.GetUpdatesChan(updateConfig)

	ctx := context.Background()
	err := t.cache.Connect(ctx)
	if err != nil {
		logger.Get().Error("Redis connection failed " + err.Error())
		panic(err)
	}

	for update := range updateChan {
		if update.Message != nil {
			chatID := update.Message.Chat.ID

			telegram_display.SetSenderChatID(chatID, t.display.(*telegram_display.BotSender))
			text := update.Message.Text

			switch text {
			case "/start":
				t.display.SendMessage("Hello! The CryptoHelper is ready for your service.")
			case "/help":
				t.display.SendMessage("CryptoHelper fetch price data from the Binance. Tap /prices to select currency and get current prices")
			case "/prices":
				if t.swapDisplay == nil {
					markupSender, err := telegram_display.NewBotMarkupSender("TELEGRAM_BOT_TOKEN")
					if err != nil {
						logger.Get().Error("NewBotMarkupSender return error " + fmt.Sprint(err))
						return err
					}
					t.swapDisplay = markupSender
				}

				t.display, t.swapDisplay = t.swapDisplay, t.display

				telegram_display.SetMarkupSenderChatID(chatID, t.display.(*telegram_display.BotMarkupSender))

				err := t.display.SendMessage("Select currency to get current price: ")

				t.display, t.swapDisplay = t.swapDisplay, t.display
				if err != nil {
					logger.Get().Error("SendMessage return error" + fmt.Sprint(err))
					return err
				}

			default:
			}
		} else {
			data := update.CallbackQuery.Data
			chatID := update.CallbackQuery.Message.Chat.ID

			var response string

			//Добавим вспомогательную функцию которая заменит куски повторяющегося кода
			getDataFunc := func(currency string) (string, error) {
				price, redisErr := t.cache.Read(ctx, currency)
				if redisErr == nil {
					logger.Get().Info(fmt.Sprintf("Redis succsecfully Read(%s) price", currency))
					response = fmt.Sprintf("%s price: "+strconv.FormatFloat(float64(price), 'f', 2, 64)+" "+currency, currency)
				} else {
					logger.Get().Error(fmt.Sprintf("Redis failed with Read %s price", currency) + redisErr.Error())
					price, err := t.GetData(currency)
					if err != nil {
						logger.Get().Error("GetData return error")
						return "", err
					}
					if redisErr == NotExist {
						err := t.cache.Write(ctx, currency, price)
						if err != nil {
							logger.Get().Error(fmt.Sprintf("Redis failed with Write %s price", currency) + redisErr.Error())
							return "", err
						}
					}
					response = fmt.Sprintf("%s price: "+strconv.FormatFloat(float64(price), 'f', 2, 64)+" USD", currency)
				}
				return response, err
			}
			switch data {
			case "btc_section":
				response, err = getDataFunc("BTC")
				if err != nil {
					logger.Get().Error("getDataFunc return error")
				}
			case "eth_section":
				response, err = getDataFunc("ETH")
				if err != nil {
					logger.Get().Error("getDataFunc return error")
				}
			case "sol_section":
				response, err = getDataFunc("SOL")
				if err != nil {
					logger.Get().Error("getDataFunc() return error")
				}
			case "ada_section":
				response, err = getDataFunc("ADA")
				if err != nil {
					logger.Get().Error("getDataFunc() return error")
				}
			case "xrp_section":
				response, err = getDataFunc("XRP")
				if err != nil {
					logger.Get().Error("getDataFunc() return error")
				}
			case "ondo_section":
				response, err = getDataFunc("ONDO")
				if err != nil {
					logger.Get().Error("getDataFunc() return error")
				}
			default:
				response = "Currency wasn't chosen"
			}

			telegram_display.SetSenderChatID(chatID, t.display.(*telegram_display.BotSender))
			t.display.SendMessage(response)

			callback := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
			t.botAPI.Request(callback)

		}

	}

	return nil
}
