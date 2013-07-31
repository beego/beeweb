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
	"strings"

	"github.com/astaxie/beego"
)

// DocsRouter serves about page.
type DocsRouter struct {
	beego.Controller
}

// Get implemented Get method for DocsRouter.
func (this *DocsRouter) Get() {
	// Set language version.
	curLang := globalSetting(this.Ctx, this.Input(), this.Data)

	this.Data["IsDocs"] = true

	reqUrl := this.Ctx.Request.URL.String()
	sec := reqUrl[strings.LastIndex(reqUrl, "/")+1:]
	if qm := strings.Index(sec, "?"); qm > -1 {
		sec = sec[:qm]
	}

	if len(sec) == 0 || sec == "docs" {
		this.Data["IsIntro"] = true
		sec = "overview"
	} else {
		this.Data["Is"+strings.Title(sec)] = true
	}

	this.Data["Title"] = strings.Title(sec)

	d, err := loadDoc(curLang.Lang + "/" + sec + ".md")
	if err != nil {
		this.Data["Data"] = err.Error()
	} else {
		this.Data["Data"] = d
	}
	this.Data["IsHasMarkdown"] = true
	this.TplNames = "docs_" + curLang.Lang + ".html"
}
