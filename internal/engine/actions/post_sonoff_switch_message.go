package actions

import (
	"bytes"
	"net/http"

	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var PostSonoffSwitchMessage types.ActionImpl = func(mm []types.Message, args types.Args, mapping types.Mapping, e types.EngineAsSupplier) (err error) {
	tpayload := types.TemplatePayload{
		Messages: mm,
	}
	areader := arguments.NewReader(&mm[0], args, mapping, &tpayload, e)
	cmd := areader.Get("Command")
	deviceId := areader.Get("DeviceId")
	err = areader.Error()
	if err != nil {
		return
	}
	device, err := e.DevicesService().GetOne(deviceId.(types.DeviceId))
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
	url := fmt.Sprintf("http://%v:%v/zeroconf/switch", ip, port)
	payload := []byte(fmt.Sprintf(`{"data":{"switch":"%v"}}`, cmd))
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err == nil && res.StatusCode == 200 {
		fmt.Println("success")
	}
	fmt.Println(url, string(payload), res, err)
	return
}
