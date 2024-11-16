package util

import (
	"encoding/json"
	"fmt"
)

func PrintAsFormattedJSON(data any) {
	if j, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Println(string(j))
	}
}
