package nullable

import (
	"encoding/json"
	"testing"
)

func TestStringUnmarshalJSON(t *testing.T) {
	var c struct {
		S String
	}
	var err error

	b := []byte(`{"S":null}`)

	if err = json.Unmarshal(b, &c); err != nil {
		t.Error(err)
	}

	if c.S != "" {
		t.Errorf("null did not unmarshal into empty string")
	}
}
