package nullable

import (
	"encoding/json"
)

type Bool bool

func (n *Bool) UnmarshalJSON(b []byte) (error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*bool)(n))
}
