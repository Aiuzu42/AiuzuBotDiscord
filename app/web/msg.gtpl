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
            <form action="/msg" method="post">
                <label for="channels">Choose a channel:</label>
                <select name="channel" id="channel">
                    {{ range $key, $value := .C }}
                        <option value="{{ $key }}">{{ $value }}</option>
                    {{end}}
                </select>
                <br>
                <br>
                Or specify channel ID:<input type="text" name="channelid">
                <br>
                <br>
                Message:
                <br>
                <textarea name="message" onkeyup="countChars(this, 2000, 'charNum');"></textarea>
                <p id="charNum">0 characters</p>
                <br>
                <br>
                <input type="submit" value="Send">
            </form>
        </div>
    </body>
</html>