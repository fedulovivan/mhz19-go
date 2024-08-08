package tbot

import (
	"log/slog"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/registry"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var withTag = logger.MakeTag("TBOT")

type service struct {
	out        engine.MessageChan
	bot        *tgbotapi.BotAPI
	botStarted bool
}

var Service engine.Service = &service{}

func (s *service) Receive() engine.MessageChan {
	return s.out
}

func (s *service) Type() engine.ChannelType {
	return engine.CHANNEL_TELEGRAM
}

func (s *service) Stop() {
	if s.botStarted {
		s.bot.StopReceivingUpdates()
	}
}

func (s *service) SendNewMessage(text string, chatID int64) {
	if chatID == 0 {
		chatID = registry.Config.TelegramChatId
	}
	msg := tgbotapi.NewMessage(chatID, text)
	// msg.ReplyToMessageID = update.Message.MessageID
	_, err := s.bot.Send(msg)
	if err != nil {
		slog.Error(withTag("Send()"), "err", err.Error())
	}
}

func (s *service) Init() {

	s.out = make(engine.MessageChan, 100)

	var err error
	s.bot, err = tgbotapi.NewBotAPI(registry.Config.TelegramToken)
	if err != nil {
		slog.Error(withTag("NewBotAPI()"), "err", err.Error())
		return
	}
	s.bot.Debug = true
	slog.Debug(withTag("Authorized"), "UserName", s.bot.Self.UserName)
	s.botStarted = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := s.bot.GetUpdatesChan(u)
	err = tgbotapi.SetLogger(slogAdapter{})
	if err != nil {
		slog.Error(withTag("SetLogger()"), "err", err.Error())
	}

	// updates
	go func() {
		for update := range updates {
			if update.Message != nil {
				// log received update
				slog.Debug(withTag("Got a message"), "UserName", update.Message.From.UserName, "Text", update.Message.Text)
				if update.Message.IsCommand() {
					slog.Debug(withTag("IsCommand() == true"), "command", update.Message.Command())
				}
				p := map[string]any{
					"Text": update.Message.Text,
					"From": update.Message.From,
				}
				outMsg := engine.Message{
					ChannelType: s.Type(),
					Timestamp:   time.Now(),
					Payload:     p,
				}
				s.out <- outMsg
				// test echo reply
				// s.SendNewMessage(update.Message.Text, update.Message.Chat.ID)
			}
		}
	}()
}
