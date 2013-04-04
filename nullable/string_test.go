package nullable

import (
	"encoding/json"
	"testing"
)

func TestStringUnmarshalJSON(t *testing.T) {
	var c struct {
		V String
	}
	var err error

	b := []byte(`{"V":null}`)

	if err = json.Unmarshal(b, &c); err != nil {
		t.Error(err)
	}

	if c.V != "" {
		t.Errorf("null did not unmarshal into empty string")
	}
}
