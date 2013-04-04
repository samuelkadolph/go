package nullable

import (
	"encoding/json"
	"testing"
)

func TestBoolUnmarshalJSON(t *testing.T) {
	var c struct {
		V Bool
	}
	var err error

	b := []byte(`{"V":null}`)

	if err = json.Unmarshal(b, &c); err != nil {
		t.Error(err)
	}

	if c.V != false {
		t.Errorf("null did not unmarshal into false")
	}
}
