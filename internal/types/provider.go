package types

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ProviderType string

const (
	PROVIDER_BURIED_DEVICES ProviderType = "buried_devices"
	PROVIDER_DNSSD          ProviderType = "dnssd"
	PROVIDER_MQTT           ProviderType = "mqtt"
	PROVIDER_SHIM_PROVIDER  ProviderType = "shim_provider"
	PROVIDER_TBOT           ProviderType = "tbot"
)

type TbotPayload struct {
	Text   string         `json:"Text"`
	From   *tgbotapi.User `json:"From"`
	ChatId int64          `json:"ChatId"`
}

type SonoffDeviceJson struct {
	ID   string `json:"Id"`
	Host string `json:"Host"`
	Ip   string `json:"Ip"`
	Port int    `json:"Port"`
	Text struct {
		ID      string `json:"id"`
		Apivers string `json:"apivers"`
		Data1   string `json:"data1"`
		Seq     string `json:"seq"`
		Txtvers string `json:"txtvers"`
		Type    string `json:"type"`
	} `json:"Text"`
}

type ZigbeeDevice struct {
	Type        string `json:"type"`
	IeeeAddress string `json:"ieee_address"`
	Definition  struct {
		Description string `json:"description"`
	} `json:"definition"`

	// Definition struct {
	// 	Exposes     []struct {
	// 		Access      int    `json:"access"`
	// 		Category    string `json:"category,omitempty"`
	// 		Description string `json:"description"`
	// 		Label       string `json:"label"`
	// 		Name        string `json:"name"`
	// 		Property    string `json:"property"`
	// 		Type        string `json:"type"`
	// 		Unit        string `json:"unit"`
	// 		ValueMax    int    `json:"value_max,omitempty"`
	// 		ValueMin    int    `json:"value_min,omitempty"`
	// 	} `json:"exposes"`
	// 	Model   string `json:"model"`
	// 	Options []struct {
	// 		Access      int    `json:"access"`
	// 		Description string `json:"description"`
	// 		Label       string `json:"label"`
	// 		Name        string `json:"name"`
	// 		Property    string `json:"property"`
	// 		Type        string `json:"type"`
	// 		ValueMax    int    `json:"value_max,omitempty"`
	// 		ValueMin    int    `json:"value_min,omitempty"`
	// 	} `json:"options"`
	// 	SupportsOta bool   `json:"supports_ota"`
	// 	Vendor      string `json:"vendor"`
	// } `json:"definition"`
	// Disabled  bool `json:"disabled"`
	// Endpoints struct {
	// 	Num1 struct {
	// 		Bindings []interface{} `json:"bindings"`
	// 		Clusters struct {
	// 			Input  []string `json:"input"`
	// 			Output []string `json:"output"`
	// 		} `json:"clusters"`
	// 		ConfiguredReportings []interface{} `json:"configured_reportings"`
	// 		Scenes               []interface{} `json:"scenes"`
	// 	} `json:"1"`
	// } `json:"endpoints"`
	// DateCode   string `json:"date_code"`
	// FriendlyName       string `json:"friendly_name"`
	// InterviewCompleted bool   `json:"interview_completed"`
	// Interviewing       bool   `json:"interviewing"`
	// Manufacturer       string `json:"manufacturer"`
	// ModelID            string `json:"model_id"`
	// NetworkAddress     int    `json:"network_address"`
	// PowerSource        string `json:"power_source"`
	// SoftwareBuildID    string `json:"software_build_id"`
	// Supported          bool   `json:"supported"`
}
