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
	"github.com/beego/beeweb/models"
)

// QuickStartRouter serves about page.
type QuickStartRouter struct {
	BaseRouter
}

// Get implemented Get method for QuickStartRouter.
func (this *QuickStartRouter) Get() {
	this.Data["IsQuickStart"] = true

	// Get language.
	curLang, _ := this.Data["LangVer"].(langType)
	df := models.GetDoc("quickstart", curLang.Lang)
	this.Data["Title"] = df.Title
	this.Data["Data"] = string(df.Data)
	this.Data["IsHasMarkdown"] = true
	this.TplNames = "quickstart_" + curLang.Lang + ".html"
}
