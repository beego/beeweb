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

// CommunityRouter serves community page.
type CommunityRouter struct {
	baseRouter
}

// Get implemented Get method for CommunityRouter.
func (this *CommunityRouter) Get() {
	this.Data["IsCommunity"] = true
	this.TplName = "community.html"

	df := models.GetDoc("usecases", this.Lang)
	this.Data["Section"] = "usecases"
	this.Data["Data"] = string(df.Data)
}
