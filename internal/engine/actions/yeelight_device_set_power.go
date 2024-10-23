package actions

import (
	"fmt"

	"github.com/SkYNewZ/go-yeelight"
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: IP, Cmd
var YeelightDeviceSetPower types.ActionImpl = func(
	compound types.MessageCompound,
	args types.Args,
	mapping types.Mapping,
	e types.EngineAsSupplier,
	tag utils.Tag,
) (err error) {
	// tpayload := types.TemplatePayload{
	// 	Messages: mm,
	// }
	// time.Sleep(time.Second)
	reader := arguments.NewReader(
		compound.Curr, args, mapping /* &tpayload */, nil, e, tag,
	)
	ip, err := arguments.GetTyped[string](&reader, "IP")
	if err != nil {
		return
	}
	cmd, err := arguments.GetTyped[string](&reader, "Cmd")
	if err != nil {
		return
	}
	y, err := yeelight.New(ip)
	if err != nil {
		return
	}
	switch cmd {
	case "On":
		err = y.On()
	case "Off":
		err = y.Off()
	default:
		err = fmt.Errorf("unsupported command '%v'", cmd)
	}
	return
}

// res, err := yeelight.Discover()
// fmt.Println("res", res)
// fmt.Println("err", err)
