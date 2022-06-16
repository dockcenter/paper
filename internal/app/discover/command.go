package discover

import (
	"encoding/json"
	"fmt"
)

func BuildCommand(workflow string, info ImageInfo) string {
	// Convert imageInfo to json
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	jsonStr := string(jsonBytes)

	return fmt.Sprintf("echo '%s' | gh workflow run %s --json", jsonStr, workflow)
}
