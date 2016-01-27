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
	"github.com/beego/beeweb/models"
)

// DonateRouter serves Donate page.
type DonateRouter struct {
	baseRouter
}

// Get implemented Get method for DonateRouter.
func (this *DonateRouter) Get() {
	this.Data["IsDonate"] = true
	this.TplName = "donate.html"

	// Get language.
	df := models.GetDoc("donate", this.Lang)
	this.Data["Title"] = df.Title
	this.Data["Data"] = string(df.Data)
	this.Data["IsHasMarkdown"] = true
}
