package push_message

import (
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

var logTag = logger.NewTag(logger.MESSAGES)

func NewApi(base *routing.RouteGroup, shimProvider types.ChannelProvider) {
	group := base.Group("/push-message")
	group.Put("", func(c *routing.Context) error {
		defer utils.TimeTrack(logTag.F, time.Now(), "api:pushMessage")
		outMsg := types.Message{
			Id:            types.MessageIdSeq.Inc(),
			Timestamp:     time.Now(),
			ChannelType:   types.CHANNEL_REST,
			DeviceClass:   types.DEVICE_CLASS_SYSTEM,
			DeviceId:      types.DEVICE_ID_FOR_THE_REST_PROVIDER_MESSAGE,
			FromEndDevice: false,
		}
		err := c.Read(&outMsg.Payload)
		if err != nil {
			return err
		}
		shimProvider.Push(outMsg)
		return c.Write(map[string]any{"ok": true})
	})
}
