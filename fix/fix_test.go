package fix

import (
	"bytes"
	"strings"
	"testing"
)

const (
	dodgy = "<dodgy><hello>Richa&rd</hello><richard.lehane@gmail.com>contracts <>f /M</dodgy>"
	clean = "<dodgy><hello>Richa&amp;rd</hello>&lt;richard.lehane@gmail.com&gt;contracts &lt;&gt;f /M</dodgy>"
)

func TestDodgy(t *testing.T) {
	rdr := strings.NewReader(dodgy)
	root := names(rdr)
	currNames := root
	var whole bool
	for _, b := range []byte("dodgy") {
		currNames, whole = currNames.next(b)
		if currNames == nil {
			t.Errorf("Dodgy fail: %v", string(b))
		}
	}
	if !whole {
		t.Error("Dodgy fail: end should be true")
	}
	currNames = root
	for _, c := range []byte("hello") {
		currNames, whole = currNames.next(c)
		if currNames == nil {
			t.Errorf("Hello fail: %v", string(c))
			break
		}
	}
	if !whole {
		t.Error("Hello fail: end should be true")
	}
}

func TestClean(t *testing.T) {
	rdr := strings.NewReader(dodgy)
	out := new(bytes.Buffer)
	err := Fix(rdr, out)
	if err != nil {
		t.Errorf("Clean fail: %v", err)
	}
	if out.String() != clean {
		t.Errorf("Clean fail: %v", out.String())
	}
}

func TestString(t *testing.T) {
	rdr := strings.NewReader(dodgy)
	root := names(rdr)
	str := root.String()
	if str != "?xml version=\"1.0\"? dodgy hello " {
		t.Errorf("String fail: expecting ?xml version=\"1.0\"? dodgy hello , got %v", str)
	}
}

func TestSanity(t *testing.T) {
	c := make(chan struct{})
	select {
	case <-c:
		t.Error("whaa")
	default:
	}
}
