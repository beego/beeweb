package models

import (
	"encoding/json"
	"os"

	"github.com/astaxie/beego"
)

type products struct {
	Projects []*Project
}

type Project struct {
	Name      string
	Thumb     string
	Desc      string
	Url       string
	Src       string
	Submitter string
	Date      string
}

var Products = new(products)

func initProuctCase() {
	fileName := "products/projects.json"

	aProducts := *Products

	var file *os.File
	var err error

	if file, err = os.Open(fileName); err != nil {
		beego.Error("open %s, %s", fileName, err.Error())
		return
	}

	d := json.NewDecoder(file)
	if err = d.Decode(&aProducts); err != nil {
		beego.Error("open %s, %s", fileName, err.Error())
		return
	}

	*Products = aProducts
}
