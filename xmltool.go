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

// Tool based on github.com/richardlehane/fixxml and github.com/richardlehane/auditxml packages.
//
// Cleans up generic XML files generated by databases. Also audits xml files, counting occurrence of elements.
// If given a directory, will do a recursive walk, fixing or auditing any files with a ".xml" extension.
//
// Examples:
//   ./xmltool -fix bad.xml > good.xml
//   ./xmltool -audit good.xml
//   ./xmltool -fix DIR_CONTAINING_BAD_XML_FILES -outdir ~/Good
//   ./xmltool -audit ~/Good -html > report.html
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/richardlehane/auditxml"
	"github.com/richardlehane/fixxml"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var fixIn = flag.String("fix", "", "fix input XML file(s)")
var auditIn = flag.String("audit", "", "perform audit function on input XML file(s)")
var html = flag.Bool("html", false, "output audit results as HTML")
var outdir = flag.String("outdir", "", "when fixing a directory, must supply path to an output directory")

var dir bool

func fix(in string, out io.Writer) error {
	inFile, err := os.Open(in)
	if err != nil {
		return err
	}
	defer inFile.Close()

	return fixxml.Fixxml(inFile, out)
}

func outpath(root string, out string, local string) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	outAbs, err := filepath.Abs(out)
	if err != nil {
		return "", err
	}
	localAbs, err := filepath.Abs(local)
	if err != nil {
		return "", err
	}
	local = strings.TrimPrefix(localAbs, rootAbs)
	return outAbs + local, nil
}

func walkFix(root string, out string) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if filepath.Ext(path) == ".xml" {
			outpath, err := outpath(root, out, path)
			if err != nil {
				return err
			}
			err = os.MkdirAll(filepath.Dir(outpath), os.ModeDir)
			if err != nil {
				return err
			}
			outFile, err := os.Create(outpath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			return fix(path, outFile)
		}
		return nil
	}
	return filepath.Walk(root, walkFn)
}

func audit(a auditxml.Audit, in string) error {
	inFile, err := os.Open(in)
	if err != nil {
		return err
	}
	defer inFile.Close()

	return a.Add(inFile)
}

func walkAudit(a auditxml.Audit, root string) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if filepath.Ext(path) == ".xml" {
			return audit(a, path)
		}
		return nil
	}
	return filepath.Walk(root, walkFn)
}

func main() {
	flag.Parse()
	if *fixIn == *auditIn {
		log.Fatal("Invalid argument: must give EITHER a -fix or -audit argument, and not both")
	}
	if len(*fixIn) > 0 {
		ffi, err := os.Stat(*fixIn)
		if err != nil {
			log.Fatal("Error opening input file %s: %v", *fixIn, err)
		}
		if ffi.IsDir() {
			if len(*outdir) < 1 {
				log.Fatal("If -fix PATH is a directory, must supply an -outdir PATH")
			}
			err = walkFix(*fixIn, *outdir)
			if err != nil {
				log.Fatalf("Error fixing xml: %v", err)
			}
		} else {
			out := new(bytes.Buffer)
			err = fix(*fixIn, out)
			if err != nil {
				log.Fatalf("Error fixing xml: %v", err)
			}
			fmt.Print(out)
		}
		return
	}
	a := make(auditxml.Audit)
	afi, err := os.Stat(*auditIn)
	if err != nil {
		log.Fatal("Error opening input file %s: %v", *auditIn, err)
	}
	if afi.IsDir() {
		err = walkAudit(a, *auditIn)
	} else {
		err = audit(a, *auditIn)
	}
	if err != nil {
		log.Fatalf("Error auditing xml: %v", err)
	}
	if *html {
		fmt.Print(a.Html())
		return
	}
	fmt.Print(a)
}
