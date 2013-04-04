package nullable

import (
	"encoding/json"
	"testing"
)

func TestInt64UnmarshalJSON(t *testing.T) {
	var c struct {
		V Int64
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
