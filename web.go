package hootail

import (
	"html/template"
	"log"
	"net/http"
)

var webPageContent = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8"/>
    <title>hootail</title>
    <script src="http://apps.bdimg.com/libs/jquery/2.1.4/jquery.min.js"></script>
    <script>

        var filterText = ""
        var logName = ""
        var ws = null

        function connect (){
            ws = new WebSocket("ws://"+ window.location.host +"/ws");
            ws.onmessage = function(e) {

                slog = eval('(' + e.data + ')');

                // 过滤查询条件
                if (filterText == "" || slog.text.indexOf(filterText) != -1){
                    $('#log').append("<pre style='color: white;font-size: 15px'>"+ slog.text +"</pre>").scrollTop($('#log')[0].scrollHeight)
                }
            };
            ws.onclose = function () {
                $('#status').css("background-color","red").text("链接断开")
                reConnect()
            }
            ws.onopen = function () {
                $('#status').css("background-color","chartreuse").text("连接成功")
                ws.send(JSON.stringify({logName:logName}));
            }
            ws.onerror = function (e) {
                $('#status').css("background-color","red").text("链接断开")
            }
        }

        function reConnect(){
            setTimeout(function(){
                connect();
            },1000);
        }
        connect();

        $(function () {
            logName = $('#logs option:first').text();

            $('#logs').on("change",function () {
                $('#log').empty()
                logName = $('#logs option:selected').text();
                ws.send(JSON.stringify({logName:logName}));
            })

            // 清屏
            $('#clear').click(function () {
                $('#log').empty()
            })

            // 过滤
            $('#filter').on('input',function () {
                filterText = $('#filter').val()
            })

        })

    </script>

</head>
<body>

<header>
    <h2 id="title">实时查看日志信息 &nbsp;<button id="status" style="background-color: darkorange">正在连接...</button>
    </h2>
    <div class="tool">

        <select id="logs">
        {{range .}}
            <option>{{ .LogName }}</option>
        {{end}}
        </select>
        <button id="clear">清屏</button>
        <span style="padding:1px;border:1px ; background:#FFF"><button style="width: auto">过滤</button><input id="filter" type="text"></span>

    </div>
</header>
<div id="log"></div>
</body>

<style>
    body {
        margin-left: 2%
    }
    #title {

    }
    #log {
        width:96%;
        height: 800px;
        background-color:#181818;
        border: 1px #ccc solid;
        overflow-y: scroll;
        margin-top: 10px;
        padding-left: 12px;
        float: left;
    }

    .tool select {
        color: blue;
        height: 30px;
        width: 120px;
        font-size: medium;
        font-weight: lighter;
    }

    .tool button {
        height: 30px;
        width: 100px;
        font-size: medium;
    }

    input {
        background-color: lightyellow;
        color: black;
        font-size: medium;
        position:absolute;
        height: 25px;
    }

</style>
</html>
`

// response page
func renderWebPage(writer http.ResponseWriter, slogs interface{}) {
	t, err := template.New("").Parse(webPageContent)
	if err != nil {
		log.Printf("renderWebPage error: %v", err)
		return
	}
	t.Execute(writer, slogs)
}
