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

package routers

import (
	"errors"
	"os"

	"github.com/astaxie/beego"
)

// QuickStartRouter serves about page.
type QuickStartRouter struct {
	beego.Controller
}

// Get implemented Get method for QuickStartRouter.
func (this *QuickStartRouter) Get() {
	// Set language version.
	curLang := globalSetting(this.Ctx, this.Input(), this.Data)

	this.Data["IsQuickStart"] = true

	d, err := loadDoc(curLang.Lang + "/quickstart.md")
	if err != nil {
		this.Data["Data"] = err.Error()
	} else {
		this.Data["Data"] = d
	}
	this.Data["IsHasMarkdown"] = true
	this.TplNames = "quickstart_" + curLang.Lang + ".html"
}

// loadDoc returns string of file data by given path.
func loadDoc(path string) (string, error) {
	f, err := os.Open("docs/" + path)
	if err != nil {
		return "", errors.New("Fail to open documentation file: " + err.Error())
	}

	fi, err := f.Stat()
	if err != nil {
		return "", errors.New("Fail to get file information: " + err.Error())
	}

	p := make([]byte, fi.Size())
	f.Read(p)
	return string(p), nil
}
