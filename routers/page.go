package routers

import (
	"github.com/beego/beeweb/models"
)

type PageRouter struct {
	baseRouter
}

func (this *PageRouter) Get() {
	this.Data["IsTeam"] = true
	this.TplName = "team.html"

	// Get language.
	df := models.GetDoc("team", this.Lang)
	this.Data["Title"] = df.Title
	this.Data["Data"] = string(df.Data)
}

type AboutRouter struct {
	baseRouter
}

func (this *AboutRouter) Get() {
	this.Data["IsAbout"] = true
	this.TplName = "about.html"

	// Get language.
	df := models.GetDoc("about", this.Lang)
	this.Data["Title"] = df.Title
	this.Data["Data"] = string(df.Data)
}
