# Beego Web

An open source project for official documentation website of beego app framework.

## Install site locally

Beego Web is a `go get` able project:

	$ go get github.com/beego/beeweb

Switch to project root path:

	$ cd $GOPATH/src/github.com/beego/beeweb

Build and run with Go tools:

	$ go build
	$ ./beeweb

Or build with bee tool:

	$ bee run beeweb

Open your browser and visit [http://localhost:8090](http://localhost:8090).

## Build as your site

This project can be easily transferred as your own documentation site, there are some tips that you may want to know:

- In the file `conf/app.conf`:

	- `langs`: languages that you want to support
	- `langNames`: user-friendly name of languages.
	- `navs`: sections' name of documentation list.
	- `samples`: samples' name of sample page.

- In the file `conf/app.ini`:
	
	- This file is ignored by rules of `.gitignore` due to it saves sensitive information; in this case, it saves Github app keys as following format:
		
			[github]
			client_id=1862bcb2******f36c
			client_secret=308d71ab53ccd858416cfceaed52******53c5f
	
	- It's not necessary to use it.

- In the file `conf/docTree.json`:

	- This file saves the file tree(with file name and commit) of your project that is hosted in Github. About how to use documentation project please see [beedoc](http://github.com/beego/beedoc). Note that if you added new section to documentation list and you do not want to wait auto-refresh, simple delete this file and restart.
	- To change the documentation project URL, you need to change it in function `checkDocUpdates` in file `models/models.go`.

- When you added new section to documentation list in `conf/app.conf`, you also need to update list in template files(`views/docs_<lang>.html`, `<lang>` depends on the language you use).

## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).