package engine_actions

import (
	"bytes"
	"log/slog"
	"net/http"

	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/fedulovivan/mhz19-go/internal/arg_reader"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var PostSonoffSwitchMessage types.ActionImpl = func(mm []types.Message, a types.Action, e types.Engine) {

	tpayload := arg_reader.TemplatePayload{
		Messages: mm,
	}
	areader := arg_reader.NewArgReader(&mm[0], a.Args, a.Mapping, &tpayload, e)
	cmd := areader.Get("Command")
	deviceId := areader.Get("DeviceId")
	if !areader.Ok() {
		slog.Error(areader.Error().Error())
		return
	}

	device, err := e.DevicesService().GetOne(deviceId.(types.DeviceId))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	djson := gabs.Wrap(device.Json)
	ip := djson.Path("ip").Data().(string)
	port := djson.Path("port").Data().(string)

	url := fmt.Sprintf("http://%v:%v/zeroconf/switch", ip, port)
	payload := []byte(fmt.Sprintf(`{"data":{"switch":"%v"}}`, cmd))
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err == nil && res.StatusCode == 200 {
		fmt.Println("success")
	}
	fmt.Println(url, string(payload), res, err)
}
