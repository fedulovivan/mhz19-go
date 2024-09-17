package utils

import (
	"encoding/json"
	"fmt"
)

func Dump(name string, value any) {
	json, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println(name, err, value)
		return
	}
	fmt.Println(name, string(json))
}
