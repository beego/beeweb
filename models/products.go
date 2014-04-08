package models

import (
	"encoding/json"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
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
	if !utils.FileExists("conf/productTree.json") {
		beego.Error("models.initBlogMap -> conf/productTree.json does not exist")
		return
	}

	f, err := os.Open("conf/productTree.json")
	if err != nil {
		beego.Error("models.initBlogMap -> load data:", err.Error())
		return
	}
	defer f.Close()

	d := json.NewDecoder(f)
	err = d.Decode(&productTree)
	if err != nil {
		beego.Error("models.initBlogMap -> decode data:", err.Error())
		return
	}

	fileName := "products/projects.json"

	aProducts := *Products

	var file *os.File

	if file, err = os.Open(fileName); err != nil {
		beego.Error("open %s, %s", fileName, err.Error())
		return
	}

	d = json.NewDecoder(file)
	if err = d.Decode(&aProducts); err != nil {
		beego.Error("open %s, %s", fileName, err.Error())
		return
	}

	for i, j := 0, len(aProducts.Projects)-1; i < j; i, j = i+1, j-1 {
		aProducts.Projects[i], aProducts.Projects[j] = aProducts.Projects[j], aProducts.Projects[i]
	}

	*Products = aProducts
}
