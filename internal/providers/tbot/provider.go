package tbot_provider

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/counters"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tag = utils.NewTag(logger.TBOT)

type provider struct {
	engine.ProviderBase
	bots    map[string]*tgbotapi.BotAPI
	botsMu  sync.RWMutex
	started chan struct{}
}

var _ types.ChannelProvider = (*provider)(nil)

func NewProvider() *provider {
	return &provider{
		bots:    make(map[string]*tgbotapi.BotAPI),
		started: make(chan struct{}),
	}
}

func (p *provider) Type() types.ProviderType {
	return types.PROVIDER_TBOT
}

func (p *provider) Channel() types.ChannelType {
	return types.CHANNEL_TELEGRAM
}

func (p *provider) Started() <-chan struct{} {
	return p.started
}

// args: botName, text, chatId
func (p *provider) Send(a ...any) error {
	botName := a[0].(string)
	chatId := a[2].(int64)
	switch text := a[1].(type) {
	case string:
		return p.SendNewMessage(text, botName, chatId)
	case []byte:
		return p.SendNewMessage(string(text), botName, chatId)
	default:
		panic(fmt.Sprintf("expected type %T", text))
	}
}

func (p *provider) Stop() {
	p.botsMu.Lock()
	defer p.botsMu.Unlock()
	if len(p.bots) == 0 {
		return
	}
	slog.Debug(tag.F("Stopping %d bot(s)...", len(p.bots)))
	for _, bot := range p.bots {
		bot.StopReceivingUpdates()
	}
	p.ProviderBase.Stop()
	// p.CloseChan()
}

func (p *provider) SendNewMessage(text string, botName string, chatId int64) (err error) {
	p.botsMu.RLock()
	defer p.botsMu.RUnlock()
	bot, found := p.bots[botName]
	if !found {
		err = fmt.Errorf("No such bot %v", botName)
		return
	}
	slog.Debug(tag.F("SendNewMessage()"), "text", text, "botName", botName, "chatId", chatId)
	msg := tgbotapi.NewMessage(chatId, text)
	_, err = bot.Send(msg)
	return
}

func (p *provider) StartBotClient(token string) (err error) {
	p.botsMu.Lock()
	defer p.botsMu.Unlock()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return
	}
	err = tgbotapi.SetLogger(slogAdapter{})
	if err != nil {
		return
	}
	bot.Debug = app.Config.TelegramDebug
	slog.Debug(tag.F("Authorized"), "UserName", bot.Self.UserName)
	p.bots[bot.Self.UserName] = bot
	updates := bot.GetUpdatesChan(
		tgbotapi.UpdateConfig{
			Offset:  0,
			Timeout: 60,
		},
	)
	go func() {
		for update := range updates {
			if update.Message != nil {
				slog.Debug(
					tag.F("Got a message"),
					"UserName", update.Message.From.UserName,
					"Text", update.Message.Text,
					"BotName", bot.Self.UserName,
				)
				if update.Message.IsCommand() {
					slog.Debug(tag.F("IsCommand() == true"), "command", update.Message.Command())
				}
				payload := types.TbotPayload{
					Text:   update.Message.Text,
					From:   update.Message.From,
					ChatId: update.Message.Chat.ID,
				}
				outMsg := types.Message{
					Id:            types.MessageIdSeq.Add(1),
					Timestamp:     time.Now(),
					DeviceId:      types.DeviceId(bot.Self.UserName),
					ChannelType:   types.CHANNEL_TELEGRAM,
					DeviceClass:   types.DEVICE_CLASS_BOT,
					Payload:       payload,
					FromEndDevice: true,
				}
				p.Push(outMsg)
			}
		}
	}()
	return
}

func (p *provider) Init() {
	p.ProviderBase.Init()
	for _, token := range app.Config.TelegramTokens {
		err := p.StartBotClient(token)
		if err != nil {
			slog.Error(tag.F("StartBotClient()"), "err", err.Error())
			counters.Inc(counters.ERRORS_ALL)
			counters.Errors.WithLabelValues(logger.MOD_TBOT).Inc()
		}
	}
	p.started <- struct{}{}
}

// chatId := app.Config.TelegramChatId
// msg.ReplyToMessageID = update.Message.MessageID
// if chatId == 0 {
// 	chatId = app.Config.TelegramChatId
// }
// chatId = int64(114333844)
