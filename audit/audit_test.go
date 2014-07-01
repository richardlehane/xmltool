package audit

import (
	"strings"
	"testing"
)

const (
	dodgy = "<dodgy><hello>Richard</hello></dodgy>"
)

func TestSingle(t *testing.T) {
	rdr := strings.NewReader(dodgy)
	a, err := Single(rdr)
	if err != nil {
		t.Errorf("Audit fail: %v", err)
	}
	d, ok := a["dodgy"]
	if !ok {
		t.Error("Audit fail: no dodgy element")
		return
	}
	if d.occurs != 1 {
		t.Errorf("Audit fail: expecting one dodgy element, got %v", d.occurs)
	}
}
