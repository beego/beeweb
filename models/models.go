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
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego"
	"github.com/slene/blackfriday"
)

var Cfg *goconfig.ConfigFile

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

var (
	docLock *sync.RWMutex
	docMap  map[string]*docFile
)

var githubCred string

func setGithubCredentials(id, secret string) {
	githubCred = "client_id=" + id + "&client_secret=" + secret
}

func init() {
	if !isExist("conf/app.ini") {
		os.Create("conf/app.ini")
	}

	var err error
	Cfg, err = goconfig.LoadConfigFile("conf/app.ini")
	if err == nil {
		beego.Info("Initialize app.ini")
	}

	setGithubCredentials(Cfg.MustValue("github", "client_id"),
		Cfg.MustValue("github", "client_secret"))

	docLock = new(sync.RWMutex)

	// Load documentation.
	initDocMap()

	// Start check ticker.
	// checkTicker = time.NewTicker(5 * time.Minute)
	// go checkTickerTimer(checkTicker.C)

	// checkDocUpdates()
}

func initDocMap() {
	// Get documentation names.
	docNames := strings.Split(beego.AppConfig.String("navs"), "|")
	docNames = append(docNames, "quickstart")
	docNames = append(docNames,
		strings.Split(beego.AppConfig.String("samples"), "|")...)

	isConfExist := isExist("conf/docTree.json")
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
			docTree.Tree = append(docTree.Tree, &docNode{Path: v})
		}
	}

	docLock.Lock()
	defer docLock.Unlock()

	docMap = make(map[string]*docFile)
	langs := []string{"zh", "en"}

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
	} else {
		return docMap[fullName]
	}
}

var checkTicker *time.Ticker

func checkTickerTimer(checkChan <-chan time.Time) {
	for {
		<-checkChan
		checkDocUpdates()
	}
}

func checkDocUpdates() {
	beego.Trace("Checking documentation updates")

	var tmpTree struct {
		Tree []*docNode
	}
	err := httpGetJSON(httpClient, "https://api.github.com/repos/beego/beedoc/git/trees/master?recursive=1&"+githubCred, &tmpTree)
	if err != nil {
		beego.Error("models.checkDocUpdates -> get trees:", err.Error())
		return
	}

	// Compare SHA.
	files := make([]*source, 0, len(tmpTree.Tree))
	for _, node := range tmpTree.Tree {
		// Skip non-md files and "README.MD".
		if !strings.HasSuffix(node.Path, ".md") || node.Path == "README.md" {
			continue
		}

		// Trim ".md".
		name := node.Path[:len(node.Path)-3]
		if checkSHA(name, node.Sha) {
			beego.Info("Need to update:", name)
			files = append(files, &source{
				name:   name,
				rawURL: "https://raw.github.com/beego/beedoc/master/" + node.Path,
			})
		}

		// For save purpose, reset name.
		node.Path = name
	}

	// Fetch files.
	if err := fetchFiles(httpClient, files, nil); err != nil {
		beego.Error("models.checkDocUpdates -> fetch files:", err.Error())
		return
	}

	// Update data.
	for _, f := range files {
		fw, err := os.Create("docs/" + f.name + ".md")
		if err != nil {
			beego.Error("models.checkDocUpdates -> open file:", err.Error())
			return
		}

		_, err = fw.Write(f.data)
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

// isExist returns if a file or directory exists
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
