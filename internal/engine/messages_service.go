package engine

import "encoding/json"

type MessagesService interface {
	Create(message Message) error
}

type messagesService struct {
	repository MessagesRepository
}

func (s messagesService) Create(message Message) (err error) {
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
