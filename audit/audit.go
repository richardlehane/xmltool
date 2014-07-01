// Copyright 2013 Richard Lehane. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package xmltool/audit reports on the contents of XML files.
//
// Example:
//   audit := audit.Single(reader)
//   fmt.Print(audit)
package audit

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strings"
)

type tag struct {
	name     string
	occurs   int
	contents int
	example  string
	files    int
}

type tags []*tag

func (t tags) Len() int           { return len(t) }
func (t tags) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t tags) Less(i, j int) bool { return t[i].contents > t[j].contents }

type Audit map[string]*tag

func (a Audit) render(f func(b *bytes.Buffer, t tags)) string {
	t := make(tags, 0, len(a))
	for _, v := range a {
		t = append(t, v)
	}
	sort.Sort(t)
	b := new(bytes.Buffer)
	f(b, t)
	return b.String()
}

// Prints the Audit
func (a Audit) String() string {
	f := func(b *bytes.Buffer, t tags) {
		for _, v := range t {
			fmt.Fprintln(b, v.name)
			fmt.Fprintf(b, "Occurs %d times in total, %d times with contents, in %d files", v.occurs, v.contents, v.files)
			if v.contents > 0 {
				fmt.Fprintf(b, "\nExample content: %s", v.example)
			}
			fmt.Fprintf(b, "\n\n")
		}
	}
	return a.render(f)
}

// Prints the Audit in a simple HTML format
func (a Audit) Html() string {
	f := func(b *bytes.Buffer, t tags) {
		fmt.Fprint(b, "<html><head><title>XML Audit</title></head><body>")
		for _, v := range t {
			fmt.Fprintf(b, "<h1>%s</h1>", v.name)
			fmt.Fprintf(b, "<p>Occurs %d times in total, %d times with contents, in %d files</p>", v.occurs, v.contents, v.files)
			if v.contents > 0 {
				fmt.Fprint(b, "<p>Example content:</p>")
				fmt.Fprintf(b, "<p>%s</p>", v.example)
			}
		}
		fmt.Fprint(b, "</body></html>")
	}
	return a.render(f)
}

// Audits a reader.
// Can be called multiple times on different readers in order to audit a set of XML files.
func (a Audit) Add(rdr io.Reader) error {
	this := make(map[string]bool)
	decoder := xml.NewDecoder(rdr)
	var err error
	var curr string
	var tg *tag
	var ok bool
	for tok, err := decoder.RawToken(); err == nil; tok, err = decoder.RawToken() {
		switch el := tok.(type) {
		case xml.StartElement:
			curr = el.Name.Local
			tg, ok = a[curr]
			if !ok {
				tg = new(tag)
				tg.name = curr
				a[curr] = tg
			}
			tg.occurs += 1
			_, ok = this[curr]
			if !ok {
				tg.files += 1
				this[curr] = true
			}
		case xml.EndElement:
			curr = ""
		case xml.CharData:
			if len(curr) < 1 {
				break
			}
			content := string(el)
			if len(strings.TrimSpace(content)) > 0 {
				tg = a[curr]
				tg.contents += 1
				if tg.contents < 2 {
					tg.example = content
				}
			}
		}
	}
	if err != io.EOF {
		return err
	}
	return nil
}

// Audit a single XML file. To Audit multiple files, make an Audit and Add readers to it.
func Single(rdr io.Reader) (Audit, error) {
	audit := make(Audit)
	err := audit.Add(rdr)
	return audit, err
}
