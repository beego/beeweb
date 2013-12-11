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
	"fmt"
	"os"
	// "strconv"
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

var docs = make(map[string]*DocRoot)

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

type oldDocNode struct {
	Sha  string
	Path string
}

// docTree descriables a documentation file structure tree.
var docTree struct {
	Tree []oldDocNode
}

var blogTree struct {
	Tree []oldDocNode
}

type docFile struct {
	Title string
	Data  []byte
}

var (
	docLock  *sync.RWMutex
	blogLock *sync.RWMutex
	docMap   map[string]*docFile
	blogMap  map[string]*docFile
)

var githubCred string

func setGithubCredentials(id, secret string) {
	githubCred = "client_id=" + id + "&client_secret=" + secret
}

func GetDocByLocale(lang string) *DocRoot {
	return docs[lang]
}

func InitModels() {
	if !com.IsFile(_CFG_PATH) {
		os.Create(_CFG_PATH)
	}

	var err error
	Cfg, err = goconfig.LoadConfigFile(_CFG_PATH)
	if err == nil {
		beego.Info("Initialize app.ini")
	} else {
		fmt.Println(err)
		os.Exit(2)
	}

	setGithubCredentials(Cfg.MustValue("github", "client_id"),
		Cfg.MustValue("github", "client_secret"))

	docLock = new(sync.RWMutex)
	blogLock = new(sync.RWMutex)

	root, err := ParseDocs("docs/zh-CN")
	if err != nil {
		beego.Error(err)
	}

	docs["zh-CN"] = root

	root, err = ParseDocs("docs/en-US")
	if err != nil {
		beego.Error(err)
	}

	docs["en-US"] = root

	// initMaps()

	// Start check ticker.
	// checkTicker = time.NewTicker(5 * time.Minute)
	// go checkTickerTimer(checkTicker.C)

	// ATTENTION: you'd better comment following code when developing.
	// if needCheckUpdate() {
	// checkFileUpdates()

	// Cfg.SetValue("app", "update_check_time", strconv.Itoa(int(time.Now().Unix())))
	// goconfig.SaveConfigFile(Cfg, _CFG_PATH)
	// }
}

func needCheckUpdate() bool {
	// Does not have record for check update.
	stamp, err := Cfg.Int64("app", "update_check_time")
	if err != nil {
		return true
	}

	if !com.IsFile("conf/docTree.json") || !com.IsFile("conf/blogTree.json") {
		return true
	}

	return time.Unix(stamp, 0).Add(5 * time.Minute).Before(time.Now())
}

func initDocMap() {
	// Load navTree.json
	fn, err := os.Open(_NAV_TREE_PATH)
	if err != nil {
		beego.Error("models.initDocMap -> load navTree.json:", err.Error())
		return
	}
	defer fn.Close()

	d := json.NewDecoder(fn)
	err = d.Decode(&navTree)
	if err != nil {
		beego.Error("models.initDocMap -> decode navTree.json:", err.Error())
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

	docNames = append(docNames, strings.Split(
		Cfg.MustValue("app", "doc_names"), "|")...)

	isConfExist := com.IsFile("conf/docTree.json")
	if isConfExist {
		f, err := os.Open("conf/docTree.json")
		if err != nil {
			beego.Error("models.initDocMap -> load data:", err.Error())
			return
		}
		defer f.Close()

		d := json.NewDecoder(f)
		err = d.Decode(&docTree)
		if err != nil {
			beego.Error("models.initDocMap -> decode data:", err.Error())
			return
		}
	} else {
		// Generate 'docTree'.
		for _, v := range docNames {
			docTree.Tree = append(docTree.Tree, oldDocNode{Path: v})
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

			docMap[fullName] = getFile("docs/" + fullName)
		}
	}
}

func initBlogMap() {
	os.Mkdir("blog", os.ModePerm)
	langs := strings.Split(Cfg.MustValue("lang", "types"), "|")
	for _, l := range langs {
		os.Mkdir("blog/"+l, os.ModePerm)
	}

	if !com.IsFile("conf/blogTree.json") {
		beego.Error("models.initBlogMap -> conf/blogTree.json does not exist")
		return
	}

	f, err := os.Open("conf/blogTree.json")
	if err != nil {
		beego.Error("models.initBlogMap -> load data:", err.Error())
		return
	}
	defer f.Close()

	d := json.NewDecoder(f)
	err = d.Decode(&blogTree)
	if err != nil {
		beego.Error("models.initBlogMap -> decode data:", err.Error())
		return
	}

	blogLock.Lock()
	defer blogLock.Unlock()

	blogMap = make(map[string]*docFile)
	for _, v := range blogTree.Tree {
		blogMap[v.Path] = getFile("blog/" + v.Path)
	}
}

func initMaps() {
	if !com.IsFile(_NAV_TREE_PATH) {
		beego.Critical(_NAV_TREE_PATH, "does not exist")
		return
	}

	initDocMap()
	initBlogMap()
}

// loadFile returns []byte of file data by given path.
func loadFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return []byte(""), errors.New("Fail to open file: " + err.Error())
	}

	fi, err := f.Stat()
	if err != nil {
		return []byte(""), errors.New("Fail to get file information: " + err.Error())
	}

	d := make([]byte, fi.Size())
	f.Read(d)
	return d, nil
}

func markdown(raw []byte) []byte {
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	htmlFlags |= blackfriday.HTML_GITHUB_BLOCKCODE
	htmlFlags |= blackfriday.HTML_OMIT_CONTENTS
	htmlFlags |= blackfriday.HTML_COMPLETE_PAGE
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_HARD_LINE_BREAK
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK

	body := blackfriday.Markdown(raw, renderer, extensions)
	return body
}

func getFile(filePath string) *docFile {
	df := &docFile{}
	p, err := loadFile(filePath + ".md")
	if err != nil {
		beego.Error("models.getFile -> ", err)
		return nil
	}

	// Parse and render.
	s := string(p)
	i := strings.Index(s, "\n")
	if i > -1 {
		// Has title.
		df.Title = strings.TrimSpace(
			strings.Replace(s[:i+1], "#", "", -1))
		df.Data = []byte(strings.TrimSpace(s[i+2:]))
	} else {
		df.Data = p
	}

	df.Data = markdown(df.Data)
	return df
}

// GetDoc returns 'docFile' by given name and language version.
func GetDoc(fullName, lang string) *docFile {
	filePath := "docs/" + lang + "/" + fullName

	if beego.RunMode == "dev" {
		return getFile(filePath)
	}

	docLock.RLock()
	defer docLock.RUnlock()
	return docMap[lang+"/"+fullName]
}

// GetBlog returns 'docFile' by given name and language version.
func GetBlog(fullName, lang string) *docFile {
	filePath := "blog/" + lang + "/" + fullName

	if beego.RunMode == "dev" {
		return getFile(filePath)
	}

	blogLock.RLock()
	defer blogLock.RUnlock()
	return blogMap[lang+"/"+fullName]
}

var checkTicker *time.Ticker

func checkTickerTimer(checkChan <-chan time.Time) {
	for {
		<-checkChan
		checkFileUpdates()
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

func checkFileUpdates() {
	beego.Trace("Checking file updates")

	type tree struct {
		ApiUrl, RawUrl, TreeName, Prefix string
	}

	var trees = []*tree{
		{
			ApiUrl:   "https://api.github.com/repos/beego/beedoc/git/trees/master?recursive=1&" + githubCred,
			RawUrl:   "https://raw.github.com/beego/beedoc/master/",
			TreeName: "conf/docTree.json",
			Prefix:   "docs/",
		},
		{
			ApiUrl:   "https://api.github.com/repos/beego/beeblog/git/trees/master?recursive=1&" + githubCred,
			RawUrl:   "https://raw.github.com/beego/beeblog/master/",
			TreeName: "conf/blogTree.json",
			Prefix:   "blog/",
		},
	}

	for _, tree := range trees {
		var tmpTree struct {
			Tree []*oldDocNode
		}

		err := com.HttpGetJSON(httpClient, tree.ApiUrl, &tmpTree)
		if err != nil {
			beego.Error("models.checkFileUpdates -> get trees:", err.Error())
			return
		}

		var saveTree struct {
			Tree []*oldDocNode
		}
		saveTree.Tree = make([]*oldDocNode, 0, len(tmpTree.Tree))

		// Compare SHA.
		files := make([]com.RawFile, 0, len(tmpTree.Tree))
		for _, node := range tmpTree.Tree {
			// Skip non-md files and "README.md".
			if !strings.HasSuffix(node.Path, ".md") || node.Path == "README.md" {
				continue
			}

			// Trim ".md".
			name := node.Path[:len(node.Path)-3]
			if checkSHA(name, node.Sha, tree.Prefix) {
				beego.Info("Need to update:", name)
				files = append(files, &rawFile{
					name:   name,
					rawURL: tree.RawUrl + node.Path,
				})
			}

			saveTree.Tree = append(saveTree.Tree, &oldDocNode{
				Path: name,
				Sha:  node.Sha,
			})
			// For save purpose, reset name.
			node.Path = name
		}

		// Fetch files.
		if err := com.FetchFiles(httpClient, files, nil); err != nil {
			beego.Error("models.checkFileUpdates -> fetch files:", err.Error())
			return
		}

		// Update data.
		for _, f := range files {
			fw, err := os.Create(tree.Prefix + f.Name() + ".md")
			if err != nil {
				beego.Error("models.checkFileUpdates -> open file:", err.Error())
				continue
			}

			_, err = fw.Write(f.Data())
			fw.Close()
			if err != nil {
				beego.Error("models.checkFileUpdates -> write data:", err.Error())
				continue
			}
		}

		// Save documentation information.
		f, err := os.Create(tree.TreeName)
		if err != nil {
			beego.Error("models.checkFileUpdates -> save data:", err.Error())
			return
		}

		e := json.NewEncoder(f)
		err = e.Encode(&saveTree)
		if err != nil {
			beego.Error("models.checkFileUpdates -> encode data:", err.Error())
			return
		}
		f.Close()
	}

	beego.Trace("Finish check file updates")
	initMaps()
}

// checkSHA returns true if the documentation file need to update.
func checkSHA(name, sha, prefix string) bool {
	var tree struct {
		Tree []oldDocNode
	}

	if prefix == "docs/" {
		tree = docTree
	} else {
		tree = blogTree
	}

	for _, v := range tree.Tree {
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
