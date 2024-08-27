package utils

import (
	"encoding/json"
	"fmt"
)

func Dump(name string, in any) {
	json, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		fmt.Println(name, err, in)
		return
	}
	fmt.Println(name, string(json))
}
