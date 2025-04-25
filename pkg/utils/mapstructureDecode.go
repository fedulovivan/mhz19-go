package utils

import "encoding/json"

// replaces mapstructure.Decode
// https://github.com/go-viper/mapstructure/issues/83
func MapstructureDecode(in any, out any) error {
	jsonData, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, out)
	if err != nil {
		return err
	}
	return nil
}
