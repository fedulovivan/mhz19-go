package messages

import (
	"encoding/json"

	// "github.com/fedulovivan/mhz19-go/internal/messages"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type MessagesService interface {
	Get() ([]types.Message, error)
	Create(message types.Message) error
}

type messagesService struct {
	repository MessagesRepository
}

func (s messagesService) Get() (messages []types.Message, err error) {
	dbmsg, err := s.repository.Get()
	if err != nil {
		return
	}
	return BuildMessages(dbmsg), nil
}

func BuildMessages(in []DbMessage) (out []types.Message) {
	for _, m := range in {
		var payload any
		_ = json.Unmarshal([]byte(m.Json), &payload)
		out = append(out, types.Message{
			ChannelType: types.ChannelType(m.ChannelTypeId),
			DeviceClass: types.DeviceClass(m.DeviceClassId),
			DeviceId:    types.DeviceId(m.DeviceId),
			Timestamp:   m.Timestamp,
			Payload:     payload,
		})
	}
	return
}

func (s messagesService) Create(message types.Message) (err error) {
	mjson, err := json.Marshal(message.Payload)
	if err != nil {
		return
	}
	_, err = s.repository.Create(DbMessage{
		ChannelTypeId: int32(message.ChannelType),
		DeviceClassId: int32(message.DeviceClass),
		DeviceId:      string(message.DeviceId),
		Timestamp:     message.Timestamp,
		Json:          string(mjson),
	})
	return
}

func NewService(r MessagesRepository) MessagesService {
	return messagesService{
		repository: r,
	}
}
