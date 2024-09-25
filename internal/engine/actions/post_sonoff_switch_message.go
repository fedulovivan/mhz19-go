package actions

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: DeviceId, Command (on,off)
var PostSonoffSwitchMessage types.ActionImpl = func(
	mm []types.Message,
	args types.Args,
	mapping types.Mapping,
	e types.EngineAsSupplier,
	tag logger.Tag,
) (err error) {
	// tpayload := types.TemplatePayload{
	// 	Messages: mm,
	// }
	reader := arguments.NewReader(
		&mm[0], args, mapping /* &tpayload */, nil, e,
	)
	deviceId, err := arguments.GetTyped[types.DeviceId](&reader, "DeviceId")
	if err != nil {
		return
	}
	command, err := arguments.GetTyped[string](&reader, "Command")
	if err != nil {
		return
	}
	device, err := e.DevicesService().GetOne(deviceId)
	if err != nil {
		return
	}
	gjson := gabs.Wrap(device.Json)
	ip, ipOk := gjson.Path("Ip").Data().(string)
	port, portOk := gjson.Path("Port").Data().(string)
	if !ipOk || !portOk {
		err = fmt.Errorf("failed to retrieve ip and port from device json: ip=%v, port=%v, %+v", ip, port, device.Json)
		return
	}
	err = httpPost(ip, port, command, tag)
	return
}

func httpPost(ip string, port string, cmd string, tag logger.Tag) error {

	url := fmt.Sprintf("http://%v:%v/zeroconf/switch", ip, port)
	payload := []byte(fmt.Sprintf(`{"data":{"switch":"%v"}}`, cmd))

	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	var bodyRaw []byte
	var bodyParsed map[string]any
	bodyRaw, err = io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyRaw, &bodyParsed)
	if err != nil {
		return err
	}

	if bodyParsed["error"] != float64(0) {
		return fmt.Errorf("%v", bodyParsed["error"])
	}

	slog.Debug(tag.F("Success"),
		"url", url,
		"request", string(payload),
		"response", string(bodyRaw),
		"status", res.StatusCode,
		"server", res.Header["Server"],
	)

	return nil
}
