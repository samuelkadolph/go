package nullable

import (
	"encoding/json"
	"testing"
)

func TestIntUnmarshalJSON(t *testing.T) {
	var c struct {
		I Int
	}
	var err error

	b := []byte(`{"I":null}`)

	if err = json.Unmarshal(b, &c); err != nil {
		t.Error(err)
	}

	if c.I != 0 {
		t.Errorf("null did not unmarshal into 0")
	}
}
