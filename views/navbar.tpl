{{define "navbar_en"}}
<div class="navbar navbar-static-top">
    <div class="navbar-inner navbar-fixed-top ">
        <div class="container">
        	<div class="brand">
	            <a class="logo" href="/">
	            	<img src="/static/img/bee.gif" style="height: 60px;">
	            </a>
				<sub class="shortintro">Simple & powerful Go web framework</sub>
        	</div>

            <ul class="nav pull-right">
                <li {{if .IsHome}}class="active"{{end}}><a href="/" title="Homepage">Home</a></li>
                <li {{if .IsAbout}}class="active"{{end}}><a href="/about" title="About Beego">About</a></li>
                <li {{if .IsCommunity}}class="active"{{end}}><a href="/community" title="The Beego Community">Community</a></li>
                <li class=""><a href="/gettingstarted" title="Getting Started">Getting started</a></li>
                <li><a href="/docs" title="Documentation">Documentation</a></li>
                <li><a target="_blank" href="http://blog.beego.me" title="Blog">Blog</a></li>
            </ul>
        </div>
    </div>
</div>
{{end}}

{{define "navbar_zh"}}
<div class="navbar navbar-static-top">
    <div class="navbar-inner navbar-fixed-top ">
        <div class="container">
            <div class="brand">
                <a class="logo" href="/">
                    <img src="/static/img/bee.gif" style="height: 60px;">
                </a>
                <sub class="shortintro">简约 & 强大并存的 Go Web 框架</sub>
            </div>

            <ul class="nav pull-right">
                <li {{if .IsHome}}class="active"{{end}}><a href="/" title="首页">首页</a></li>
                <li {{if .IsAbout}}class="active"{{end}}><a href="/about" title="关于 Beego">关于 Beego</a></li>
                <li {{if .IsCommunity}}class="active"{{end}}><a href="/community" title="开发者社区">开发者社区</a></li>
                <li class=""><a href="/gettingstarted" title="快速入门">快速入门</a></li>
                <li><a href="/docs" title="API 文档">API 文档</a></li>
                <li><a target="_blank" href="http://blog.beego.me" title="官方博客">博客</a></li>
            </ul>
        </div>
    </div>
</div>
{{end}}