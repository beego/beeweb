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
	"github.com/beego/beeweb/models"
)

// SamplesRouter serves about page.
type SamplesRouter struct {
	beego.Controller
}

// Get implemented Get method for SamplesRouter.
func (this *SamplesRouter) Get() {
	// Set language version.
	curLang := globalSetting(this.Ctx, this.Input(), this.Data)

	this.Data["IsSamples"] = true

	reqUrl := this.Ctx.Request.URL.String()
	sec := reqUrl[strings.LastIndex(reqUrl, "/")+1:]
	if qm := strings.Index(sec, "?"); qm > -1 {
		sec = sec[:qm]
	}

	if len(sec) == 0 || sec == "samples" {
		this.Redirect("/samples/Samples_Introduction", 302)
		return
	} else {
		this.Data[sec] = true
	}

	df := models.GetDoc(sec, curLang.Lang)
	if df == nil {
		this.Redirect("/samples/Samples_Introduction", 302)
		return
	}

	this.Data["Title"] = df.Title
	this.Data["Data"] = string(df.Data)
	this.Data["IsHasMarkdown"] = true
	this.TplNames = "samples_" + curLang.Lang + ".html"
}
