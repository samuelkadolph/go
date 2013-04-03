package nullable

import (
	"encoding/json"
	"testing"
)

func TestInt64UnmarshalJSON(t *testing.T) {
	var c struct {
		I Int64
	}
	var err error

	b := []byte(`{"b":null}`)

	if err = json.Unmarshal(b, &c); err != nil {
		t.Error(err)
	}

	if c.I != 0 {
		t.Errorf("null did not unmarshal into 0")
	}
}
