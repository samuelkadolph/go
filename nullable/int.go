package nullable

import (
	"encoding/json"
)

type Int int

func (n *Int) UnmarshalJSON(b []byte) (error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*int)(n))
}
