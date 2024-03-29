<html>
    <head>
    <title></title>
    </head>
    <script>
        function countChars(obj, max, elem){
            var maxLength = max;
            var strLength = obj.value.length;
            if(strLength > maxLength){
                document.getElementById(elem).innerHTML = '<span style="color: red;">'+strLength+' de '+maxLength+' caracteres</span>';
            }else{
                document.getElementById(elem).innerHTML = strLength+' de '+maxLength+' caracteres';
            }
        }
    </script>
    <body>
        <div>
            <form action="/msgembed" method="post">
                <label for="channels">Choose a channel:</label>
                <select name="channel" id="channel">
                    {{ range $key, $value := .C }}
                        <option value="{{ $key }}">{{ $value }}</option>
                    {{end}}
                </select>
                <br>
                <br>
                Or specify channel ID:<input type="text" name="channelid" value="{{ .Emb.ChannelID}}">
                <br>
                <br>
                <label for="colorpicker">Color Picker:</label>
                <input type="color" id="colorpicker" name="color" value={{ .HexColor }}>
                <br>
                <br>
                Titulo:<textArea name="title" onkeyup="countChars(this, 256, 'titlecharnum');">{{ .Emb.Title }}</textArea>
                <br>
				<p id="titlecharnum">0 characters</p>
                <br>
                <br>
                Descripción:<textArea name="description" onkeyup="countChars(this, 2048, 'desccharnum');">{{ .Emb.Content }}</textArea>
                <br>
				<p id="desccharnum">0 characters</p>
                <span style="color: red;">{{ .ErrorMsg }}</span>
                <br>
                <br>
                Imagen:<input type="text" name="image" value="{{ .Emb.Image}}">
                <br>
                <br>
                Thumbnail:<input type="text" name="thumbnail" value="{{ .Emb.Thumbnail}}">
                <br>
                <br>
                Delay in minutes:<input type="text" name="delay" value="0">
                <br>
                <br>
                <input type="submit" value="Send">
                <br>
                <br>
                Fields
                <br>
                <br>
                {{ range $i, $element := .Emb.Fields }}
				Name {{ $i }}:<textarea name="field{{ $i }}title" onkeyup="countChars(this, 256, 'field{{ $i }}titlecharnum');">{{ $element.Name }}</textarea>
				<br>
				<p id="field{{ $i }}titlecharnum">0 characters</p>
				Value {{ $i }}:<textarea name="field{{ $i }}content" onkeyup="countChars(this, 1024, 'field{{ $i }}contentcharnum');">{{ $element.Value }}</textarea>
				<br>
				<p id="field{{ $i }}contentcharnum">0 characters</p>
				Inline {{ $i }}:<input type="checkbox" name="field{{ $i }}inline" value="true" {{ if $element.Inline }} checked {{ end }}>
				<br>
				<br>
                {{ end }}
            </form>
        </div>
    </body>
</html>