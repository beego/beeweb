{{define "navbar"}}
<div class="navbar navbar-static-top">
    <div class="navbar-inner navbar-fixed-top ">
        <div class="container">
        	<div class="brand">
	            <a class="logo" href="/">
	            	<img src="/static/img/bee.gif" style="height: 60px;">
	            </a>
				<sub class="shortintro">{{i18n .Lang "app_intro"}}</sub>
        	</div>

            <ul class="nav pull-right">
                <li {{if .IsHome}}class="active"{{end}}><a href="/">{{i18n .Lang "home"}}</a></li>
                <li {{if .IsAbout}}class="active"{{end}}><a href="/about">{{i18n .Lang "about"}}</a></li>
                <li {{if .IsQuickStart}}class="active"{{end}}><a href="/quickstart">{{i18n .Lang "getting started"}}</a></li>
                <li {{if .IsCommunity}}class="active"{{end}}><a href="/community">{{i18n .Lang "community"}}</a></li>
                <li {{if .IsDocs}}class="active"{{end}}><a href="/docs">{{i18n .Lang "docs"}}</a></li>
                <li><a target="_blank" href="http://blog.beego.me">{{i18n .Lang "blog"}}</a></li>
            </ul>
        </div>
    </div>
</div>
{{end}}