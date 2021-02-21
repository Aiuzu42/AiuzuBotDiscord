<html>
    <head>
    <title></title>
    </head>
    <body>
        <div>
            <form action="/edit" method="post">
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
                Message ID: <input type="text" name="messageid">
                <br>
                <br>
                <input type="submit" value="Send">
            </form>
        </div>
    </body>
</html>