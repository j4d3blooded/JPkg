package jpkg_bin

import (
	"bytes"
	"testing"
)

func TestReadBinary(t *testing.T) {

	b := bytes.Buffer{}

	if err := BinaryWrite(&b, test{"XYZ", 128, true}); err != nil {
		t.Logf("error binary writing: %v", err.Error())
		t.FailNow()
	}

	v, err := BinaryRead[test](&b)
	if err != nil {
		t.Logf("error binary reading: %v", err.Error())
		t.FailNow()
	}

	if v.Name != "XYZ" {
		t.Logf("name incorrectly serialized: %v", v.Name)
		t.FailNow()
	}

	if v.Value != 128 {
		t.Logf("value incorrectly serialized: %v", v.Value)
		t.FailNow()
	}

	if v.Boolean != true {
		t.Logf("boolean incorrectly serialized: %v", v.Boolean)
		t.FailNow()
	}
}

type test struct {
	Name    string
	Value   uint64
	Boolean bool
}
