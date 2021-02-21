<html>
    <head>
    <title></title>
    </head>
    <script>
        function countChars(obj){
            var maxLength = 2000;
            var strLength = obj.value.length;
            if(strLength > maxLength){
                document.getElementById("charNum").innerHTML = '<span style="color: red;">'+strLength+' de '+maxLength+' caracteres</span>';
            }else{
                document.getElementById("charNum").innerHTML = strLength+' de '+maxLength+' caracteres';
            }
        }
</script>
    <body>
        <div>
            <form action="/editmsg" method="post">
                Channel ID: <input type="text" name="channelid" readonly value="{{ .ChannelID }}">
                <br>
                <br>
                Message ID: <input type="text" name="messageid" readonly value="{{ .MessageID }}">
                <br>
                <br>
                Message:
                <br>
                <textarea name="message" onkeyup="countChars(this);">{{ .Content }}</textarea>
                <p id="charNum">{{ .ContentL }} characters</p>
                <br>
                <br>
                <input type="submit" value="Send">
            </form>
        </div>
    </body>
</html>