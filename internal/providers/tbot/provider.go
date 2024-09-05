package tbot_provider

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var logTag = logger.MakeTag(logger.TBOT)

type provider struct {
	engine.ProviderBase
	bot        *tgbotapi.BotAPI
	botStarted bool
}

var Provider types.ChannelProvider = &provider{}

func (p *provider) Channel() types.ChannelType {
	return types.CHANNEL_TELEGRAM
}

func (p *provider) Send(a ...any) {
	first := a[0]
	switch fTyped := first.(type) {
	case string:
		p.SendNewMessage(fTyped, 0)
	case []byte:
		p.SendNewMessage(string(fTyped), 0)
	default:
		panic(fmt.Sprintf("expected type %T", fTyped))
	}
}

func (p *provider) Stop() {
	slog.Debug(logTag("Stopping bot..."))
	if p.botStarted {
		p.bot.StopReceivingUpdates()
	} else {
		slog.Warn(logTag("Not started"))
	}
}

func (p *provider) SendNewMessage(text string, chatId int64) {
	// msg.ReplyToMessageID = update.Message.MessageID
	slog.Debug(logTag("SendNewMessage()"), "text", text, "chatId", chatId)
	if chatId == 0 {
		chatId = app.Config.TelegramChatId
	}
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := p.bot.Send(msg)
	if err != nil {
		slog.Error(logTag("SendNewMessage()"), "err", err.Error())
	}
}

func (p *provider) Init() {

	p.Out = make(types.MessageChan, 100)

	var err error
	p.bot, err = tgbotapi.NewBotAPI(app.Config.TelegramToken)
	if err != nil {
		slog.Error(logTag("NewBotAPI()"), "err", err.Error())
		return
	}
	p.bot.Debug = app.Config.TelegramDebug
	slog.Debug(logTag("Authorized"), "UserName", p.bot.Self.UserName)
	p.botStarted = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := p.bot.GetUpdatesChan(u)
	err = tgbotapi.SetLogger(slogAdapter{})
	if err != nil {
		slog.Error(logTag("SetLogger()"), "err", err.Error())
	}

	// updates
	go func() {
		for update := range updates {
			if update.Message != nil {
				slog.Debug(logTag("Got a message"), "UserName", update.Message.From.UserName, "Text", update.Message.Text)
				if update.Message.IsCommand() {
					slog.Debug(logTag("IsCommand() == true"), "command", update.Message.Command())
				}
				payload := map[string]any{
					"Text":   update.Message.Text,
					"From":   update.Message.From,
					"ChatId": update.Message.Chat.ID,
				}
				outMsg := types.Message{
					DeviceId:    types.DeviceId(p.bot.Self.UserName),
					ChannelType: types.CHANNEL_TELEGRAM,
					DeviceClass: types.DEVICE_CLASS_BOT,
					Timestamp:   time.Now(),
					Payload:     payload,
				}
				p.Out <- outMsg
			}
		}
	}()
}
