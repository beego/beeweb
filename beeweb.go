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

// An open source project for official documentation website of beego app framework.
package main

import (
	"os"

	"github.com/astaxie/beego"
	"github.com/beego/beeweb/models"
	"github.com/beego/beeweb/routers"
)

const (
	APP_VER = "0.5.2.0924"
)

// We have to call a initialize function manully
// because we use `bee bale` to pack static resources
// and we cannot make sure that which init() execute first.
func initialize() {
	models.InitModels()
	routers.InitRouter()

	// Set App version and log level.
	beego.AppName = models.Cfg.MustValue("beego", "app_name")
	beego.RunMode = models.Cfg.MustValue("beego", "run_mode")
	beego.HttpPort = models.Cfg.MustInt("beego", "http_port_"+beego.RunMode)

	routers.IsPro = beego.RunMode == "pro"
	if routers.IsPro {
		beego.SetLevel(beego.LevelInfo)
		os.Mkdir("./log", os.ModePerm)
		beego.BeeLogger.SetLogger("file", `{"filename": "log/log"}`)
	} else {
		// beewatch.Start(beewatch.Trace)r
	}
}

func main() {
	initialize()

	beego.Info(beego.AppName, APP_VER)

	// Register routers.
	beego.Router("/", &routers.HomeRouter{})
	beego.Router("/about", &routers.AboutRouter{})
	beego.Router("/community", &routers.CommunityRouter{})
	beego.Router("/quickstart", &routers.QuickStartRouter{})
	beego.Router("/docs", &routers.DocsRouter{})
	beego.Router("/docs/:all", &routers.DocsRouter{})
	beego.Router("/samples", &routers.SamplesRouter{})
	beego.Router("/samples/:all", &routers.SamplesRouter{})

	// Register template functions.

	beego.Run()
}
