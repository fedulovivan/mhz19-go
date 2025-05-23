package actions

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: DeviceId, Command (on,off)
var PostSonoffSwitchMessage types.ActionImpl = func(
	compound types.MessageCompound,
	args types.Args,
	mapping types.Mapping,
	e types.ServiceAndProviderSupplier,
	tag utils.Tag,
) (err error) {
	reader := arguments.NewReader(
		compound.Curr, args, mapping, nil, e, tag,
	)
	deviceId, err := arguments.GetTyped[types.DeviceId](&reader, "DeviceId")
	if err != nil {
		return
	}
	command, err := arguments.GetTyped[string](&reader, "Command")
	if err != nil {
		return
	}
	device, err := e.GetDevicesService().GetOne(deviceId)
	if err != nil {
		return
	}
	out := new(types.SonoffDeviceJson)
	err = utils.MapstructureDecode(device.Json, out)
	if err != nil {
		return
	}
	if out.Ip == "" || out.Port == 0 {
		err = fmt.Errorf("failed to retrieve Ip or Port from device json: Ip=%v, Port=%d, %+v", out.Ip, out.Port, device.Json)
		return
	}
	err = httpPost(out.Ip, out.Port, command, tag)
	return
}

func httpPost(ip string, port int, cmd string, tag utils.Tag) error {

	url := fmt.Sprintf("http://%v:%d/zeroconf/switch", ip, port)
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
