package messages

import (
	"database/sql"
	"encoding/json"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var _ types.MessagesService = (*messagesService)(nil)

type messagesService struct {
	repository MessagesRepository
}

func (s messagesService) Get() (messages []types.Message, err error) {
	dbmsg, err := s.repository.Get(sql.NullString{})
	if err != nil {
		return
	}
	return BuildMessages(dbmsg), nil
}

func (s messagesService) GetByDeviceId(deviceId string) (messages []types.Message, err error) {
	dbmsg, err := s.repository.Get(db.NewNullString(deviceId))
	if err != nil {
		return
	}
	return BuildMessages(dbmsg), nil
}

func BuildMessages(in []DbMessage) (out []types.Message) {
	for _, m := range in {
		var payload map[string]any
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

func NewService(r MessagesRepository) types.MessagesService {
	return messagesService{
		repository: r,
	}
}
