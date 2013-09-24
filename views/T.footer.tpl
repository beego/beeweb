{{define "footer_en"}}
<div class="footer-landscape">
    <div class="footer-landscape-image">
        <!-- footer -->
        <div class="container">
            <div class="row">
                <div class="span12 footer">
                    <div class="span8 tbox">
                    	Beego App framework is an open source project, launched by <a href="https://github.com/astaxie">AstaXie</a>, and maintained by beego developer community. Under the <a href="http://www.apache.org/licenses/LICENSE-2.0.html">apache 2.0 licence</a>.
                        <strong>Language:</strong>
					    <div class="btn-group dropup">
						    <button class="btn dropdown-toggle" data-toggle="dropdown">{{.CurLang}} <span class="caret"></span></button>
						    <ul class="dropdown-menu">
							{{$keyword := .Keyword}}
						    	{{range .RestLangs}}
						    	<li><a href="?lang={{.Lang}}&q={{$keyword}}">{{.Name}}</a></li>
						    	{{end}}
						    </ul>
					    </div>
                    </div>

                    <div class="span4 tbox textright social links">
                        <a class="twitter" href="https://twitter.com/xiemengjun">Twitter</a>
                        <a class="github" href="https://github.com/astaxie/beego">GitHub</a>
                        <a class="googleplus" href="https://plus.google.com/111292884696033638814">Goolge+</a>
                    </div>
                </div>
            </div>
        </div>
        <!-- end of footer -->
    </div>
</div>
{{template "static_file" .}}
{{end}}

{{define "footer_zh"}}
<div class="footer-landscape">
    <div class="footer-landscape-image">
        <!-- footer -->
        <div class="container">
            <div class="row">
                <div class="span12 footer">
                    <div class="span8 tbox">
                    	beego 应用框架是一个开源项目，初经 <a href="https://github.com/astaxie">Asta谢</a> 发起，现由 beego 开发者社区维护。授权许可遵循 <a href="http://www.apache.org/licenses/LICENSE-2.0.html">apache 2.0 licence</a>。
                        <strong>语言选项：</strong>
					    <div class="btn-group dropup">
						    <button class="btn dropdown-toggle" data-toggle="dropdown">{{.CurLang}} <span class="caret"></span></button>
						    <ul class="dropdown-menu">
							{{$keyword := .Keyword}}
						    	{{range .RestLangs}}
						    	<li><a href="?lang={{.Lang}}&q={{$keyword}}">{{.Name}}</a></li>
						    	{{end}}
						    </ul>
					    </div>
                    </div>

                    <div class="span4 tbox textright social links">
                        <a class="twitter" href="https://twitter.com/xiemengjun">Twitter</a>
                        <a class="github" href="https://github.com/astaxie/beego">GitHub</a>
                        <a class="googleplus" href="https://plus.google.com/111292884696033638814">Goolge+</a>
                    </div>
                </div>
            </div>
        </div>
        <!-- end of footer -->
    </div>
</div>
{{template "static_file" .}}
{{end}}

{{define "static_file"}}
<script type="text/javascript" src="http://ajax.googleapis.com/ajax/libs/jquery/1.10.1/jquery.min.js"></script>
<script src="/static/js/bootstrap.min.js"></script>
<script src="/static/js/jquery.scrollTo-min.js"></script>
<script src="/static/js/jquery-scrolltofixed-min.js"></script>
   
{{if .IsHome}}
<script type="text/javascript">
    function moveRight() {
        $('#tweets-container').scrollTo( '+=620px', 300 );
    }
    function moveLeft() {
        $('#tweets-container').scrollTo( '-=620px', 300 );
    }
    function showLeftShadow() {
        $('.tweet-navigator-left').addClass('shadow').show();
    }

    $('document').ready(function() {
        $('.tweet-navigator-right').click(moveRight);
        $('.tweet-navigator-left').click(moveLeft);

        $('#tweets-container').scroll(showLeftShadow);
    })
</script>
{{end}}

{{if .IsHasMarkdown}}
<script type="text/javascript" src="/static/js/highlight.pack.js"></script>
<script type="text/javascript">
    (function($){

        var doc = $("#markdown");

        // highlight
    	doc.find("code").each(function(k, e){
    		var code = $(e);
    		if (code.attr("class") == undefined) {
    			code.attr("class", "no-highlight");
    		}
    		hljs.highlightBlock(code.get(0), hljs.tabReplace);
    	});

        // Encode url.
        doc.find("a").each(function(){
            var node = $(this);
            var link = node.attr("href");
            var index = link.indexOf("#");
            if (link.indexOf("http") == 0 && link.indexOf(window.location.hostname) == -1) {
                return
            }
            if (index <  0 || index+1 > link.length) {
                return;
            }
            var val = link.substring(index+1, link.length);
            val = encodeURIComponent(decodeURIComponent(val).toLowerCase().replace(/\s+/g, "-"));
            var end = index;
            if (index-3 > 0 && link.substring(index-3, index) == ".md") {
                end = index - 3;
            };
            node.attr("href", link.substring(0, end)+"#"+val);
        });

        // Set anchor.
        doc.find("h1, h2, h3, h4, h5, h6").each(function(){
            var node = $(this);
            var val = encodeURIComponent(node.text().toLowerCase().replace(/\s+/g, "-"));
            node = node.wrap('<div id="' + val + '" class="anchor-wrap" ></div>');
            node.append('<a class="anchor" href="#' + val + '"><span class="octicon octicon-link"></span></a>')
        });

    })(jQuery);
$(function(){
        $('#navlist').scrollToFixed({
            marginTop: $('.container').outerHeight(true) + 20,
            
            zIndex: 999
        });

        if (($("#navlist").height()+80)>$(window).height()){
                $("#navlist").css({"overflow-y":"scroll","height":($(window).height()-100)})
            }else{
                $("#navlist").css({"overflow-y":"","height":"auto"})
            }

        $(window).resize(function () {
            if (($("#navlist").height()+80)>$(window).height()){
                $("#navlist").css({"overflow-y":"scroll","height":($(window).height()-100)})
            }else{
                $("#navlist").css({"overflow-y":"","height":"auto"})
            }
        });
        
        
    });
</script>
{{end}}

{{if .IsPro}}
<script>
    (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
    (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
    m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
    })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

    ga('create', 'UA-40109089-3', 'beego.me');
    ga('send', 'pageview');
</script>
{{end}}
{{end}}
