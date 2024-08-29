package engine

import (
	"encoding/json"

	"github.com/fedulovivan/mhz19-go/internal/messages"
)

type MessagesService interface {
	Get() ([]Message, error)
	Create(message Message) error
}

type messagesService struct {
	repository messages.MessagesRepository
}

func (s messagesService) Get() (messages []Message, err error) {
	dbmsg, err := s.repository.Get()
	if err != nil {
		return
	}
	return BuildMessages(dbmsg), nil
}

func BuildMessages(in []messages.DbMessage) (out []Message) {
	for _, m := range in {
		var payload any
		_ = json.Unmarshal([]byte(m.Json), &payload)
		out = append(out, Message{
			ChannelType: ChannelType(m.ChannelTypeId),
			DeviceClass: DeviceClass(m.DeviceClassId),
			DeviceId:    DeviceId(m.DeviceId),
			Timestamp:   m.Timestamp,
			Payload:     payload,
		})
	}
	return
}

func (s messagesService) Create(message Message) (err error) {
	mjson, err := json.Marshal(message.Payload)
	if err != nil {
		return
	}
	_, err = s.repository.Create(messages.DbMessage{
		ChannelTypeId: int32(message.ChannelType),
		DeviceClassId: int32(message.DeviceClass),
		DeviceId:      string(message.DeviceId),
		Timestamp:     message.Timestamp,
		Json:          string(mjson),
	})
	return
}

func NewMessagesService(r messages.MessagesRepository) MessagesService {
	return messagesService{
		repository: r,
	}
}
