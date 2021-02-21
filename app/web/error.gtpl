<html>
    <head>
    <title></title>
    </head>
    <body>
        <div>
            <h1>{{ .Code }}</h1>
            <form action="/" method="get">
                {{ .Message }}
                <br>
                <br>
                <input type="submit" value="Return">
            </form>
        </div>
    </body>
</html>