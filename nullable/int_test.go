package nullable

import (
	"encoding/json"
	"testing"
)

func TestIntUnmarshalJSON(t *testing.T) {
	var c struct {
		V Int
	}
	var err error

	b := []byte(`{"V":null}`)

	if err = json.Unmarshal(b, &c); err != nil {
		t.Error(err)
	}

	if c.V != 0 {
		t.Errorf("null did not unmarshal into 0")
	}
}
