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

package routers

import (
	"strings"

	"github.com/beego/beeweb/models"
)

// BlogRouter serves about page.
type BlogRouter struct {
	baseRouter
}

// Get implemented Get method for BlogRouter.
func (this *BlogRouter) Get() {
	this.Data["IsBlog"] = true
	this.TplName = "blog.html"

	reqUrl := this.Ctx.Request.URL.String()
	fullName := reqUrl[strings.LastIndex(reqUrl, "/")+1:]
	if qm := strings.Index(fullName, "?"); qm > -1 {
		fullName = fullName[:qm]
	}

	df := models.GetBlog(fullName, this.Lang)
	if df == nil {
		this.Redirect("/blog", 302)
		return
	}

	this.Data["Title"] = df.Title
	this.Data["Data"] = string(df.Data)
	this.Data["IsHasMarkdown"] = true
}
