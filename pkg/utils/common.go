package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJson(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("unmarshal error: %v", err.Error())
	}
	return string(b)
}
