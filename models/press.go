// Copyright 2013 Beego Web authors
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

package models

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type DocList []*DocNode

func (s DocList) Len() int           { return len(s) }
func (s DocList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s DocList) Less(i, j int) bool { return s[i].Sort < s[j].Sort }

type DocNode struct {
	root        bool
	IsDir       bool
	Path        string
	RelPath     string
	FileRelPath string
	FilePath    string
	Date        time.Time
	Name        string
	Sort        int
	Link        string
	Docs        DocList
	dirs        map[string]*DocNode
	Root        *DocRoot
	Parent      *DocNode
}

func (d *DocNode) SortDocs() {
	sort.Sort(d.Docs)
}

func (d *DocNode) HasContent() bool {
	return len(d.FilePath) > 0
}

func (d *DocNode) GetContent() string {
	if !d.HasContent() {
		return ""
	}

	body, err := ioutil.ReadFile(d.FilePath)
	if err != nil {
		return ""
	}

	if i := bytes.Index(body, []byte("---")); i != -1 {
		body = body[i+3:]
		if i = bytes.Index(body, []byte("---")); i != -1 {
			body = body[i+3:]
			i = 0
			m := 0
		mFor:
			for {
				if len(body) > 0 {
					if body[0] == ' ' || body[0] == '\n' {
						if body[0] == '\n' {
							m += 1
						}
						if m == 2 {
							break mFor
						}
					} else {
						break mFor
					}
					body = body[1:]
				} else {
					break mFor
				}
			}

			return string(markdown(body))
		}
	}

	return ""
}

type DocRoot struct {
	Wd    string
	Path  string
	Doc   *DocNode
	links map[string]*DocNode
}

func (d *DocRoot) GetNodeByLink(link string) (*DocNode, bool) {
	n, ok := d.links[link]
	return n, ok
}

func (d *DocRoot) walkParse() error {
	var err error
	if d.Path, err = filepath.Abs(d.Path); err != nil {
		return err
	}

	defer func() {
		if err == nil {
			d.sortAll(d.Doc)
		}
	}()

	err = filepath.Walk(d.Path, d.walk)
	return err
}

func (d *DocRoot) sortAll(node *DocNode) {
	for _, n := range node.Docs {
		if n.IsDir {
			d.sortAll(n)
		}
	}
	node.SortDocs()
}

func (d *DocRoot) makeDirNode(path string) error {
	relPath, _ := filepath.Rel(d.Path, path)

	var docDir *DocNode

	if d.Doc == nil {
		d.Doc = new(DocNode)
		d.Doc.dirs = make(map[string]*DocNode)
		docDir = d.Doc

	} else {
		list := strings.Split(relPath, string(filepath.Separator))
		node := d.Doc
		for _, p := range list {
			if n, ok := node.dirs[p]; ok {
				node = n
			} else {
				n = new(DocNode)
				n.dirs = make(map[string]*DocNode)
				n.Parent = node
				node.Docs = append(node.Docs, n)
				node.dirs[p] = n
				node = n
			}
		}

		docDir = node
	}

	docDir.Root = d
	docDir.Path = path
	docDir.RelPath = relPath
	docDir.IsDir = true

	return nil
}

func (d *DocRoot) getDirNode(path string) *DocNode {
	node := d.Doc
	list := strings.Split(path, string(filepath.Separator))
	for _, p := range list {
		if n, ok := node.dirs[p]; ok {
			node = n
		}
	}
	return node
}

func (d *DocRoot) makeFileNode(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	relPath, _ := filepath.Rel(d.Path, path)
	relPath = strings.Replace(relPath, "\\", "/", -1)

	docDir := d.getDirNode(filepath.Dir(relPath))

	var bingo bool
	var doc *DocNode
	rd := bufio.NewReader(file)
	no := 0
	for {
		line, _, err := rd.ReadLine()
		if err == io.EOF {
			break
		}

		if no > 3 && !bingo {
			break
		}

		if no > 20 && bingo {
			return fmt.Errorf("document %s not contained ended tag `---`", path)
		}

		data := string(bytes.TrimSpace(line))

		if len(data) == 3 && data == "---" {

			if bingo {
				if doc.root {
					if len(docDir.FilePath) > 0 {
						return fmt.Errorf("node %s has a document %s, can not replicate by %s",
							docDir.Path, docDir.FilePath, path)
					}

					docDir.Name = doc.Name
					docDir.Date = doc.Date
					docDir.Link = doc.Link
					docDir.Sort = doc.Sort

				mFor:
					for {
						l, _, er := rd.ReadLine()
						if er != nil {
							break mFor
						}
						if len(bytes.TrimSpace(l)) > 0 {
							docDir.FilePath = path
							break mFor
						}
					}

					if len(docDir.Link) == 0 {
						docDir.Link = docDir.RelPath + "/"
					}

					docDir.FileRelPath = relPath

					doc = docDir
				} else {
					doc.RelPath = relPath
					doc.FilePath = path
					if len(doc.Link) == 0 {
						doc.Link = doc.RelPath
						// doc.Link = strings.TrimSuffix(doc.RelPath, filepath.Ext(doc.RelPath))
					}

					docDir.Docs = append(docDir.Docs, doc)
				}

				if dc, ok := d.links[doc.Link]; ok {
					return fmt.Errorf("document %s's link %s is already used by %s", path, doc.Link, dc.Path)
				}

				d.links[doc.Link] = doc

				break
			}

			doc = new(DocNode)
			doc.Path = path
			doc.Root = d
			doc.Parent = docDir

			bingo = true
		}

		if bingo {
			parts := strings.SplitN(data, ":", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				switch name {
				case "root":
					doc.root, _ = strconv.ParseBool(value)
				case "name":
					doc.Name = value
				case "date":
					doc.Date, err = beego.DateParse(value, "Y-m-d H:i")
					if err != nil {
						return err
					}
				case "link":
					doc.Link = value
				case "sort":
					n, _ := strconv.ParseInt(value, 10, 64)
					doc.Sort = int(n)
				}
			}
		}
	}

	return nil
}

func (d *DocRoot) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return filepath.SkipDir
	}

	if !info.IsDir() && info.Size() == 0 {
		return nil
	}

	if info.IsDir() {
		if err := d.makeDirNode(path); err != nil {
			return err
		}
	} else {
		return d.makeFileNode(path)
	}

	return nil
}

func ParseDocs(path string) (*DocRoot, error) {
	root := new(DocRoot)
	root.Path = path
	root.links = make(map[string]*DocNode)

	if err := root.walkParse(); err == nil {
		return root, err
	} else {
		return nil, err
	}
}
