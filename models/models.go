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
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"github.com/astaxie/beego/utils"
	"github.com/slene/blackfriday"
)

var docs = make(map[string]*DocRoot)

type oldDocNode struct {
	Sha  string
	Path string
	Type string
}

// docTree descriables a documentation file structure tree.
var docTree struct {
	Tree []oldDocNode
}

var blogTree struct {
	Tree []oldDocNode
}

var productTree struct {
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

	setGithubCredentials(beego.AppConfig.String("github::client_id"),
		beego.AppConfig.String("github::client_secret"))

	docLock = new(sync.RWMutex)
	blogLock = new(sync.RWMutex)

	parseDocs()
	initMaps()
	initProuctCase()

	updateTask := toolbox.NewTask("check file update", "0 */5 * * * *", checkFileUpdates)

	if needCheckUpdate() {
		if err := updateTask.Run(); err != nil {
			beego.Error(err)
		}

		beego.AppConfig.Set("app::update_check_time", strconv.Itoa(int(time.Now().Unix())))
	}

	// ATTENTION: you'd better comment following code when developing.
	toolbox.AddTask("check file update", updateTask)
	toolbox.StartTask()
}

func parseDocs() {
	root, err := ParseDocs("docs/zh-CN")
	if err != nil {
		beego.Error(err)
	}

	if root != nil {
		docs["zh-CN"] = root
	}

	root, err = ParseDocs("docs/en-US")
	if err != nil {
		beego.Error(err)
	}

	if root != nil {
		docs["en-US"] = root
	}
}

func needCheckUpdate() bool {
	// Does not have record for check update.
	stamp, err := beego.AppConfig.Int64("app::update_check_time")
	if err != nil {
		return true
	}

	if !utils.FileExists("conf/docTree.json") || !utils.FileExists("conf/blogTree.json") ||
		!utils.FileExists("conf/productTree.json") {
		return true
	}

	return time.Unix(stamp, 0).Add(5 * time.Minute).Before(time.Now())
}

func initDocMap() {
	// Documentation names.
	docNames := make([]string, 0, 20)
	docNames = append(docNames, strings.Split(
		beego.AppConfig.String("app::doc_names"), "|")...)

	isConfExist := utils.FileExists("conf/docTree.json")
	if isConfExist {
		f, err := os.Open("conf/docTree.json")
		if err != nil {
			beego.Error("models.initDocMap -> load data:", err.Error())
			return
		}
		defer f.Close()

		d := json.NewDecoder(f)
		if err = d.Decode(&docTree); err != nil {
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
	langs := strings.Split(beego.AppConfig.String("lang::types"), "|")

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
	langs := strings.Split(beego.AppConfig.String("lang::types"), "|")
	for _, l := range langs {
		os.Mkdir("blog/"+l, os.ModePerm)
	}

	if !utils.FileExists("conf/blogTree.json") {
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
	initDocMap()
	initBlogMap()
	initProuctCase()
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
	if strings.Contains(filePath, "images") {
		return nil
	}

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
		if len(s) >= i+2 {
			df.Data = []byte(strings.TrimSpace(s[i+2:]))
		}
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

func checkFileUpdates() error {
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
		{
			ApiUrl:   "https://api.github.com/repos/beego/products/git/trees/master?recursive=1&" + githubCred,
			RawUrl:   "https://raw.github.com/beego/products/master/",
			TreeName: "conf/productTree.json",
			Prefix:   "products/",
		},
	}

	for _, tree := range trees {
		var tmpTree struct {
			Tree []*oldDocNode
		}

		err := getHttpJson(tree.ApiUrl, &tmpTree)
		if err != nil {
			return errors.New("models.checkFileUpdates -> get trees: " + err.Error())
		}

		var saveTree struct {
			Tree []*oldDocNode
		}
		saveTree.Tree = make([]*oldDocNode, 0, len(tmpTree.Tree))

		// Compare SHA.
		files := make([]*rawFile, 0, len(tmpTree.Tree))
		for _, node := range tmpTree.Tree {
			// Skip non-md files and "README.md".
			if node.Type != "blob" || (!strings.HasSuffix(node.Path, ".md") &&
				!strings.Contains(node.Path, "images") &&
				!strings.HasSuffix(node.Path, ".json")) ||
				strings.HasPrefix(strings.ToLower(node.Path), "readme") {
				continue
			}

			name := strings.TrimSuffix(node.Path, ".md")

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
		if err := getFiles(files); err != nil {
			return errors.New("models.checkFileUpdates -> fetch files: " + err.Error())
		}

		// Update data.
		for _, f := range files {
			os.MkdirAll(path.Join(tree.Prefix, path.Dir(f.name)), os.ModePerm)
			suf := ".md"
			if strings.Contains(f.name, "images") ||
				strings.HasSuffix(f.name, ".json") {
				suf = ""
			}
			fw, err := os.Create(tree.Prefix + f.name + suf)
			if err != nil {
				beego.Error("models.checkFileUpdates -> open file:", err.Error())
				continue
			}

			_, err = fw.Write(f.data)
			fw.Close()
			if err != nil {
				beego.Error("models.checkFileUpdates -> write data:", err.Error())
				continue
			}
		}

		// Save documentation information.
		f, err := os.Create(tree.TreeName)
		if err != nil {
			return errors.New("models.checkFileUpdates -> save data: " + err.Error())
		}

		e := json.NewEncoder(f)
		err = e.Encode(&saveTree)
		if err != nil {
			return errors.New("models.checkFileUpdates -> encode data: " + err.Error())
		}
		f.Close()
	}

	beego.Trace("Finish check file updates")
	parseDocs()
	initMaps()
	return nil
}

// checkSHA returns true if the documentation file need to update.
func checkSHA(name, sha, prefix string) bool {
	var tree struct {
		Tree []oldDocNode
	}

	switch prefix {
	case "docs/":
		tree = docTree
	case "blog/":
		tree = blogTree
	default:
		tree = productTree
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
