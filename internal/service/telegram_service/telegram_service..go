package telegram_service

import (
	"context"
	cacheinterface "cryptoHelper/internal/cache/cache_interface"
	"cryptoHelper/internal/datasource/datasource_interface"
	"cryptoHelper/internal/display/display_interface"
	logger "cryptoHelper/pkg/applogger"
	"errors"
	"fmt"
	"strconv"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

const NotExist = redis.Nil

type Telegram struct {
	botAPI        *tgBotAPI.BotAPI
	datasource    datasource_interface.Datasource
	markupDisplay display_interface.Display
	swapDisplay   display_interface.Display
	cache         cacheinterface.CacheHandler
}

func NewTelegram(token string) (*Telegram, error) {
	bot, err := tgBotAPI.NewBotAPI(token)
	if err != nil {
		logger.Get().Error("Telegram wasn't created NewBotAPI return err")
		return nil, err
	} else {
		logger.Get().Debug("Telegram successfully created botAPI")
		return &Telegram{botAPI: bot}, nil
	}
}

func (t *Telegram) Message(chatID int64, message string) tgBotAPI.MessageConfig {
	msg := tgBotAPI.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"

	return msg
}

func (t *Telegram) MarkupMessage(chatID int64, message string) tgBotAPI.MessageConfig {
	msg := tgBotAPI.NewMessage(chatID, message)

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
	return msg
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

// Этот костыль возник из за ненужности display внутри Telegram
// и желания соответсвовать архитектуре Display, Datasource, Service
// Эта функция нужна чтобы реализовывать интерфейс Service
func (t *Telegram) SetOutput(dspl display_interface.Display) error {
	return nil
}

func (t Telegram) GetData(currencyName any) (float64, error) {
	if currencyName, ok := currencyName.(string); ok {
		if t.datasource != nil {
			logger.Get().Debug("Datasource_handler called ExtractCurrentPrice() successfully")
			return t.datasource.ExtractCurrentPrice(currencyName)
		} else {
			logger.Get().Error("Datasource_interface is nil, GetData() operation cannot be completed")
			return 0.0, errors.New("datasource_interface is nil, operation can not be completed")
		}
	} else {
		logger.Get().Error("GetData::Cannot convert currencyName any into string")
		return 0.0, errors.New("GetData::Cannot convert currencyName any into string")
	}
}

func (t Telegram) SendData(message any) error {
	if message, ok := message.(tgBotAPI.MessageConfig); ok {
		_, err := t.botAPI.Send(message)
		if err != nil {
			logger.Get().Error("Send() return error in SendMessage func")
			return err
		} else {
			logger.Get().Debug("Send() successfully send message")
			return nil
		}
	} else {
		logger.Get().Error("SendData::Cannot convert massage any into tgBotAPI.MessageConfig")
		return errors.New("SendData::Cannot convert massage any into tgBotAPI.MessageConfig")
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

	errChan := make(chan (error))

	for {
		select {
		case update := <-updateChan:

			if update.Message != nil {
				go func(update tgBotAPI.Update) {
					chatID := update.Message.Chat.ID

					text := update.Message.Text

					switch text {
					case "/start":
						msg := t.Message(chatID, "Hello! The CryptoHelper is ready for your service.")
						err := t.SendData(msg)
						if err != nil {
							logger.Get().Error("SendData in /text section return error")
							errChan <- err
						}
					case "/help":
						msg := t.Message(chatID,
							"CryptoHelper fetch price data from the Binance. Tap /prices to select currency and get current prices")
						err := t.SendData(msg)
						if err != nil {
							logger.Get().Error("SendData in /help section return error")
							errChan <- err
						}
					case "/prices":

						msg := t.MarkupMessage(chatID, "Select currency to get current price: ")
						err := t.SendData(msg)
						if err != nil {
							logger.Get().Error("SendData in /prices section return error")
							errChan <- err
						}

					default:
					}
				}(update)
			} else {
				go func(ctx context.Context, update tgBotAPI.Update) {

					chatID := update.CallbackQuery.Message.Chat.ID
					data := update.CallbackQuery.Data
					var response string

					//Добавим вспомогательную функцию которая заменит куски повторяющегося кода
					getDataFunc := func(currency string) (string, error) {
						price, redisErr := t.cache.Read(ctx, currency)
						if redisErr == nil {
							logger.Get().Info(fmt.Sprintf("Redis succsecfully Read(%s) price", currency))
							response = fmt.Sprintf("%s price: "+strconv.FormatFloat(float64(price), 'f', 2, 64)+" USD", currency)
						} else {

							price, err := t.GetData(currency)
							if err != nil {
								logger.Get().Error("GetData return error")
								return "", err
							}
							if redisErr == NotExist {
								logger.Get().Debug(fmt.Sprintf("Redis hasn't got %s price  ", currency) + redisErr.Error())
								err := t.cache.Write(ctx, currency, price)
								if err != nil {
									logger.Get().Error(fmt.Sprintf("Redis failed with Write %s price  ", currency) + redisErr.Error())
									return "", err
								}
							} else {
								logger.Get().Error(fmt.Sprintf("Redis failed with Read %s price  ", currency) + redisErr.Error())
								return "", redisErr
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
							errChan <- err
						}
					case "eth_section":
						response, err = getDataFunc("ETH")
						if err != nil {
							logger.Get().Error("getDataFunc return error")
							errChan <- err
						}
					case "sol_section":
						response, err = getDataFunc("SOL")
						if err != nil {
							logger.Get().Error("getDataFunc() return error")
							errChan <- err
						}
					case "ada_section":
						response, err = getDataFunc("ADA")
						if err != nil {
							logger.Get().Error("getDataFunc() return error")
							errChan <- err
						}
					case "xrp_section":
						response, err = getDataFunc("XRP")
						if err != nil {
							logger.Get().Error("getDataFunc() return error")
							errChan <- err
						}
					case "ondo_section":
						response, err = getDataFunc("ONDO")
						if err != nil {
							logger.Get().Error("getDataFunc() return error")
							errChan <- err
						}
					default:
						response = "Currency wasn't chosen"
					}

					msg := t.MarkupMessage(chatID, response)
					err := t.SendData(msg)
					if err != nil {
						logger.Get().Error("SendData return error while sending response")
						errChan <- err
					}

					callback := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
					t.botAPI.Request(callback)
				}(ctx, update)
			}
		case err := <-errChan:
			logger.Get().Error("Update return err " + err.Error())
		}
	}

}
