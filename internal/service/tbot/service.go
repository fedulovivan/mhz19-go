package tbot_service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var logTag = logger.MakeTag("TBOT")

type service struct {
	engine.ServiceBase
	bot        *tgbotapi.BotAPI
	botStarted bool
}

var Service engine.Service = &service{}

func (s *service) Channel() engine.ChannelType {
	return engine.CHANNEL_TELEGRAM
}

func (s *service) Send(a ...any) {
	first := a[0]
	switch fTyped := first.(type) {
	case string:
		s.SendNewMessage(fTyped, 0)
	case []byte:
		s.SendNewMessage(string(fTyped), 0)
	default:
		panic(fmt.Sprintf("expected type %T", fTyped))
	}
}

func (s *service) Stop() {
	if s.botStarted {
		slog.Debug(logTag("Stopping bot..."))
		s.bot.StopReceivingUpdates()
	}
}

func (s *service) SendNewMessage(text string, chatId int64) {
	// msg.ReplyToMessageID = update.Message.MessageID
	slog.Debug(logTag("SendNewMessage()"), "text", text, "chatId", chatId)
	if chatId == 0 {
		chatId = app.Config.TelegramChatId
	}
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := s.bot.Send(msg)
	if err != nil {
		slog.Error(logTag("SendNewMessage()"), "err", err.Error())
	}
}

func (s *service) Init() {

	s.Out = make(engine.MessageChan, 100)

	var err error
	s.bot, err = tgbotapi.NewBotAPI(app.Config.TelegramToken)
	if err != nil {
		slog.Error(logTag("NewBotAPI()"), "err", err.Error())
		return
	}
	s.bot.Debug = false
	// s.bot.Debug = app.Config.TelegramDebug
	slog.Debug(logTag("Authorized"), "UserName", s.bot.Self.UserName)
	s.botStarted = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := s.bot.GetUpdatesChan(u)
	err = tgbotapi.SetLogger(slogAdapter{})
	if err != nil {
		slog.Error(logTag("SetLogger()"), "err", err.Error())
	}

	// updates
	go func() {
		for update := range updates {
			if update.Message != nil {

				// log received update
				// fmt.Println(update.Message.Chat.ID)

				slog.Debug(logTag("Got a message"), "UserName", update.Message.From.UserName, "Text", update.Message.Text)
				if update.Message.IsCommand() {
					slog.Debug(logTag("IsCommand() == true"), "command", update.Message.Command())
				}
				p := map[string]any{
					// "Command": update.Message.Command(),
					"Text":   update.Message.Text,
					"From":   update.Message.From,
					"ChatId": update.Message.Chat.ID,
				}
				outMsg := engine.Message{
					ChannelType: engine.CHANNEL_TELEGRAM,
					DeviceClass: engine.DEVICE_CLASS_BOT,
					Timestamp:   time.Now(),
					Payload:     p,
				}
				s.Out <- outMsg
				// test echo reply
				// s.SendNewMessage(update.Message.Text, update.Message.Chat.ID)
			}
		}
	}()
}
