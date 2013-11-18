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

// Package models is for loading and updating documentation files.
package models

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/com"
	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego"
	"github.com/slene/blackfriday"
)

const (
	_CFG_PATH      = "conf/app.ini"
	_NAV_TREE_PATH = "conf/navTree.json"
)

var Cfg *goconfig.ConfigFile

type navNode struct {
	Name  string
	Nodes []string
}

// navTree descriables the navigation structure tree.
var navTree struct {
	Tree []navNode
}

type node struct {
	Index          int
	Name, FullName string
}

type Section struct {
	Name  string
	Nodes []node
}

var TplTree struct {
	Sections []Section
}

type docNode struct {
	Sha  string
	Path string
}

// docTree descriables a documentation file structure tree.
var docTree struct {
	Tree []docNode
}

type docFile struct {
	Title string
	Data  []byte
}

var (
	docLock *sync.RWMutex
	docMap  map[string]*docFile
)

var githubCred string

func setGithubCredentials(id, secret string) {
	githubCred = "client_id=" + id + "&client_secret=" + secret
}

func InitModels() {
	if !com.IsFile(_CFG_PATH) {
		os.Create(_CFG_PATH)
	}

	var err error
	Cfg, err = goconfig.LoadConfigFile(_CFG_PATH)
	if err == nil {
		beego.Info("Initialize app.ini")
	}

	setGithubCredentials(Cfg.MustValue("github", "client_id"),
		Cfg.MustValue("github", "client_secret"))

	docLock = new(sync.RWMutex)

	// Load documentation.
	initDocMap()

	// ATTENTION: you'd better comment following code when developing.
	if beego.RunMode == "pro" {
		// Start check ticker.
		checkTicker = time.NewTicker(5 * time.Minute)
		go checkTickerTimer(checkTicker.C)

		checkDocUpdates()
	}
}

func initDocMap() {
	if !com.IsFile(_NAV_TREE_PATH) {
		beego.Critical(_NAV_TREE_PATH, "does not exist")
		return
	}

	// Load navTree.json
	fn, err := os.Open(_NAV_TREE_PATH)
	if err != nil {
		beego.Error("models.init -> load navTree.json:", err.Error())
		return
	}
	defer fn.Close()

	d := json.NewDecoder(fn)
	err = d.Decode(&navTree)
	if err != nil {
		beego.Error("models.init -> decode navTree.json:", err.Error())
		return
	}

	// Documentation names.
	docNames := make([]string, 0, 20)

	// Generate usable TplTree for template.
	TplTree.Sections = make([]Section, len(navTree.Tree))
	for i, sec := range navTree.Tree {
		TplTree.Sections[i].Name = sec.Name
		TplTree.Sections[i].Nodes = make([]node, len(sec.Nodes))
		for j, nod := range sec.Nodes {
			TplTree.Sections[i].Nodes[j].Index = j + 1
			TplTree.Sections[i].Nodes[j].Name = nod

			docName := sec.Name + "_" + nod
			TplTree.Sections[i].Nodes[j].FullName = docName
			docNames = append(docNames, docName)
		}
	}

	docNames = append(docNames, "quickstart")
	docNames = append(docNames, "donate")
	docNames = append(docNames, "usecases")

	isConfExist := com.IsFile("conf/docTree.json")
	if isConfExist {
		f, err := os.Open("conf/docTree.json")
		if err != nil {
			beego.Error("models.init -> load data:", err.Error())
			return
		}
		defer f.Close()

		d := json.NewDecoder(f)
		err = d.Decode(&docTree)
		if err != nil {
			beego.Error("models.init -> decode data:", err.Error())
			return
		}
	} else {
		// Generate 'docTree'.
		for _, v := range docNames {
			docTree.Tree = append(docTree.Tree, docNode{Path: v})
		}
	}

	docLock.Lock()
	defer docLock.Unlock()

	docMap = make(map[string]*docFile)
	langs := strings.Split(Cfg.MustValue("lang", "types"), "|")

	os.Mkdir("docs", os.ModePerm)
	for _, l := range langs {
		os.Mkdir("docs/"+l, os.ModePerm)
		for _, v := range docTree.Tree {
			var fullName string
			if isConfExist {
				fullName = v.Path
			} else {
				fullName = l + "/" + v.Path
			}

			docMap[fullName] = getDoc(fullName)
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

func markdown(raw []byte) (out []byte) {
	return blackfriday.MarkdownCommon(raw)
}

func getDoc(fullName string) *docFile {
	df := &docFile{}
	d, err := loadDoc(fullName + ".md")
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

		df.Data = markdown(df.Data)
	}

	return df
}

// GetDoc returns 'docFile' by given name and language version.
func GetDoc(path, lang string) *docFile {
	docLock.RLock()
	defer docLock.RUnlock()

	fullName := lang + "/" + path

	if beego.RunMode == "dev" {
		return getDoc(fullName)
	}
	return docMap[fullName]
}

var checkTicker *time.Ticker

func checkTickerTimer(checkChan <-chan time.Time) {
	for {
		<-checkChan
		checkDocUpdates()
	}
}

type rawFile struct {
	name   string
	rawURL string
	data   []byte
}

func (rf *rawFile) Name() string {
	return rf.name
}

func (rf *rawFile) RawUrl() string {
	return rf.rawURL
}

func (rf *rawFile) Data() []byte {
	return rf.data
}

func (rf *rawFile) SetData(p []byte) {
	rf.data = p
}

func checkDocUpdates() {
	beego.Trace("Checking documentation updates")

	var tmpTree struct {
		Tree []*docNode
	}
	err := com.HttpGetJSON(httpClient, "https://api.github.com/repos/beego/beedoc/git/trees/master?recursive=1&"+githubCred, &tmpTree)
	if err != nil {
		beego.Error("models.checkDocUpdates -> get trees:", err.Error())
		return
	}

	// Compare SHA.
	files := make([]com.RawFile, 0, len(tmpTree.Tree))
	for _, node := range tmpTree.Tree {
		// Skip non-md files and "README.MD".
		if !strings.HasSuffix(node.Path, ".md") || node.Path == "README.md" {
			continue
		}

		// Trim ".md".
		name := node.Path[:len(node.Path)-3]
		if checkSHA(name, node.Sha) {
			beego.Info("Need to update:", name)
			files = append(files, &rawFile{
				name:   name,
				rawURL: "https://raw.github.com/beego/beedoc/master/" + node.Path,
			})
		}

		// For save purpose, reset name.
		node.Path = name
	}

	// Fetch files.
	if err := com.FetchFiles(httpClient, files, nil); err != nil {
		beego.Error("models.checkDocUpdates -> fetch files:", err.Error())
		return
	}

	// Update data.
	for _, f := range files {
		fw, err := os.Create("docs/" + f.Name() + ".md")
		if err != nil {
			beego.Error("models.checkDocUpdates -> open file:", err.Error())
			return
		}

		_, err = fw.Write(f.Data())
		fw.Close()
		if err != nil {
			beego.Error("models.checkDocUpdates -> write data:", err.Error())
			return
		}
	}

	beego.Trace("Finish check documentation updates")
	initDocMap()

	// Save documentation information.
	f, err := os.Create("conf/docTree.json")
	if err != nil {
		beego.Error("models.checkDocUpdates -> save data:", err.Error())
		return
	}
	defer f.Close()

	e := json.NewEncoder(f)
	err = e.Encode(&tmpTree)
	if err != nil {
		beego.Error("models.checkDocUpdates -> encode data:", err.Error())
		return
	}
}

// checkSHA returns true if the documentation file need to update.
func checkSHA(name, sha string) bool {
	for _, v := range docTree.Tree {
		if v.Path == name {
			// Found.
			if v.Sha != sha {
				// Need to update.
				return true
			}
			return false
		}
	}
	// Not found.
	return true
}
