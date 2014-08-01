package audit

import (
	"strings"
	"testing"
)

const (
	dodgy = "<dodgy><hello>Richard</hello></dodgy>"
)

func TestXMLAuditSingle(t *testing.T) {
	rdr := strings.NewReader(dodgy)
	a, err := XMLAuditSingle(rdr)
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

func TestSingleTagAudit(t *testing.T) {
	rdr := strings.NewReader(dodgy)
	ta, err := TagAuditSingle(rdr, "hello")
	if err != nil {
		t.Errorf("Tag audit fail: %v", err)
	}
	contents := (*ta)[0].contents
	if len(contents) < 1 {
		t.Error("Tag audit fail: no hello contents")
		return
	}
	if contents[0] != "Richard" {
		t.Errorf("Tag audit fail: expecting 'Richard' contents, got %v", contents[0])
	}
}
