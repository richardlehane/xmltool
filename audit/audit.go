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

type Audit interface {
	Add(io.Reader) error
	Html() string
	String() string
}

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

type XMLAudit map[string]*tag

func (a XMLAudit) render(f func(b *bytes.Buffer, t tags)) string {
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
func (a XMLAudit) String() string {
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
func (a XMLAudit) Html() string {
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
func (a XMLAudit) Add(rdr io.Reader) error {
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
func XMLAuditSingle(rdr io.Reader) (XMLAudit, error) {
	audit := make(XMLAudit)
	err := audit.Add(rdr)
	return audit, err
}

type tagContents struct {
	name     string
	contents []string
}

// TagAudit records the contents of particular tags you are interested in inspecting.
// It is very simplistic and doesn't allow any nesting, so don't try to view the contents of tags if they have hierarchical relations with each other!
type TagAudit []*tagContents

func NewTagAudit(tags ...string) *TagAudit {
	tg := make(TagAudit, len(tags))
	for i, nm := range tags {
		tg[i] = new(tagContents)
		tg[i].name = nm
		tg[i].contents = make([]string, 0, 100)
	}
	return &tg
}

func (t *TagAudit) checkNm(nm string) (int, bool) {
	for i, ta := range *t {
		if ta.name == nm {
			return i, true
		}
	}
	return 0, false
}

// Tag Audits a reader.
// Can be called multiple times on different readers in order to tag audit a set of XML files.
func (t *TagAudit) Add(rdr io.Reader) error {
	decoder := xml.NewDecoder(rdr)
	var err error
	var curr string
	var ok bool
	for tok, err := decoder.RawToken(); err == nil; tok, err = decoder.RawToken() {
		switch el := tok.(type) {
		case xml.StartElement:
			if _, ok = t.checkNm(el.Name.Local); ok {
				curr = el.Name.Local
			}
		case xml.EndElement:
			if curr == el.Name.Local {
				curr = ""
			}
		case xml.CharData:
			if len(curr) < 1 {
				break
			}
			content := string(el)
			if len(strings.TrimSpace(content)) > 0 {
				idx, _ := t.checkNm(curr)
				(*t)[idx].contents = append((*t)[idx].contents, content)
			}
		}
	}
	if err != io.EOF {
		return err
	}
	return nil
}

// Prints the Tag Audit
func (t *TagAudit) String() string {
	b := new(bytes.Buffer)
	for _, v := range *t {
		fmt.Fprintln(b, v.name)
		fmt.Fprintln(b, "Contents:")
		for _, c := range v.contents {
			fmt.Fprintln(b, c)
		}
		fmt.Fprintln(b, "")
	}
	return b.String()
}

// Prints the Tag Audit in a simple HTML format
func (t *TagAudit) Html() string {
	b := new(bytes.Buffer)
	fmt.Fprint(b, "<html><head><title>XML Audit</title></head><body>")
	for _, v := range *t {
		fmt.Fprintf(b, "<h1>%s</h1>", v.name)
		fmt.Fprintln(b, "<h2>Contents</h2>")
		for _, c := range v.contents {
			fmt.Fprintf(b, "<p>%s</p>", c)
		}
	}
	fmt.Fprint(b, "</body></html>")
	return b.String()
}

// Tag Audit a single XML file. For multiple files, create a new TagAudit and Add readers to it.
func TagAuditSingle(rdr io.Reader, tags ...string) (*TagAudit, error) {
	ta := NewTagAudit(tags...)
	err := ta.Add(rdr)
	return ta, err
}
