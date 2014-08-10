package mockerfile

import (
	"encoding/json"
	"strings"
)

// Base on JSON but accept comment which has prefix '#'
func Unmarshal(data []byte, v interface{}) error {
	edited := []byte{}
	for _, s := range strings.Split(string(data), "\n") {
		if trimmed := strings.TrimLeft(s, " "); !strings.HasPrefix(trimmed, "#") {
			edited = append(edited, []byte(trimmed)...)
		}
	}
	return json.Unmarshal(edited, v)
}
