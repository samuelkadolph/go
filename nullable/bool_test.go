package nullable

import (
	"encoding/json"
	"testing"
)

func TestBoolUnmarshalJSON(t *testing.T) {
	var c struct {
		B Bool
	}
	var err error

	b := []byte(`{"B":null}`)

	if err = json.Unmarshal(b, &c); err != nil {
		t.Error(err)
	}

	if c.B != false {
		t.Errorf("null did not unmarshal into false")
	}
}
