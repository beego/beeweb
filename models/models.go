// Copyright 2013 Unknown
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package models is for loading and updating documentation files.
package models

import (
	"errors"
	"os"
	"strings"

	"github.com/astaxie/beego"
)

type docNode struct {
	Sha  string
	Path string
}

// docTree descriables a documentation file strcuture tree.
var docTree struct {
	Tree []*docNode
}

type docFile struct {
	Title string
	Data  []byte
}

var docMap map[string]*docFile

func init() {
	// Get documentation names.
	docNames := strings.Split(beego.AppConfig.String("navs"), "|")
	docNames = append(docNames, "quickstart")

	// Generate 'docTree'.
	for _, v := range docNames {
		docTree.Tree = append(docTree.Tree, &docNode{Path: v})
	}

	// Load documentation.
	docMap = make(map[string]*docFile)
	langs := []string{"zh", "en"}
	for _, l := range langs {
		for _, v := range docTree.Tree {
			d, err := loadDoc(l + "/" + v.Path + ".md")
			df := &docFile{}
			if err != nil {
				df.Data = []byte(err.Error())
			} else {
				s := string(d)
				i := strings.Index(s, "\n")
				if i > -1 {
					// Has title.
					df.Title = strings.TrimSpace(
						strings.Replace(s[:i+1], "#", "", -1))
					df.Data = []byte(strings.TrimSpace(s[i+2:]))
				} else {
					df.Data = d
				}
			}

			docMap[l+"_"+v.Path] = df
		}
	}
}

// loadDoc returns []byte of file data by given path.
func loadDoc(path string) ([]byte, error) {
	f, err := os.Open("docs/" + path)
	if err != nil {
		return []byte(""), errors.New("Fail to open documentation file: " + err.Error())
	}

	fi, err := f.Stat()
	if err != nil {
		return []byte(""), errors.New("Fail to get file information: " + err.Error())
	}

	d := make([]byte, fi.Size())
	f.Read(d)
	return d, nil
}

// GetDoc returns 'docFile' by given name and language version.
func GetDoc(path, lang string) *docFile {
	return docMap[lang+"_"+path]
}
