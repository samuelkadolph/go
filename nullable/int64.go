package nullable

import (
	"encoding/json"
)

type Int64 int64

func (n *Int64) UnmarshalJSON(b []byte) (error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*int64)(n))
}
