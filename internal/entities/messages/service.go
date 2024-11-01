package messages

import (
	"database/sql"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/goccy/go-json"
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

func (s messagesService) GetWithTemperature(deviceId types.DeviceId) ([]types.TemperatureMessage, error) {
	messages, err := s.repository.GetWithTemperature(
		db.NewNullString(string(deviceId)),
	)
	if err != nil {
		return nil, err
	}
	result := make([]types.TemperatureMessage, 0, len(messages))
	for _, m := range messages {
		r := types.TemperatureMessage{
			Timestamp: time.Unix(int64(m.Timestamp), 0),
		}
		if m.Temperature.Valid {
			r.Temperature = m.Temperature.Float64
		}
		result = append(result, r)
	}
	return result, nil
}

func (s messagesService) GetByDeviceId(deviceId types.DeviceId) (messages []types.Message, err error) {
	dbmsg, err := s.repository.Get(db.NewNullString(string(deviceId)))
	if err != nil {
		return
	}
	return BuildMessages(dbmsg), nil
}

func BuildMessages(in []DbMessage) []types.Message {
	out := make([]types.Message, 0, len(in))
	for _, m := range in {
		var payload map[string]any
		_ = json.Unmarshal([]byte(m.Json), &payload)
		out = append(out, types.Message{
			FromEndDevice: true,
			ChannelType:   types.ChannelType(m.ChannelTypeId),
			DeviceClass:   types.DeviceClass(m.DeviceClassId),
			DeviceId:      types.DeviceId(m.DeviceId),
			Timestamp:     m.Timestamp,
			Payload:       payload,
		})
	}
	return out
}

func ToDb(message types.Message) (res DbMessage, err error) {
	mjson, err := json.Marshal(message.Payload)
	if err != nil {
		return
	}
	res = DbMessage{
		ChannelTypeId: int32(message.ChannelType),
		DeviceClassId: int32(message.DeviceClass),
		DeviceId:      string(message.DeviceId),
		Timestamp:     message.Timestamp,
		Json:          string(mjson),
	}
	return
}

func (s messagesService) CreateAll(messages []types.Message) error {
	res := make([]DbMessage, 0, len(messages))
	for _, message := range messages {
		dbMessage, err := ToDb(message)
		if err != nil {
			return err
		}
		res = append(res, dbMessage)
	}
	return s.repository.CreateAll(res)
}

func (s messagesService) Create(message types.Message) (err error) {
	dbMessage, err := ToDb(message)
	if err != nil {
		return
	}
	err = s.repository.CreateAll([]DbMessage{dbMessage})
	return
}

func NewService(r MessagesRepository) types.MessagesService {
	return messagesService{
		repository: r,
	}
}
