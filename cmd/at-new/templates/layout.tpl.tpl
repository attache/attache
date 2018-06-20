<!DOCTYPE html>
<html>
    <head>
        <title>{{block "title" .}}[[.Name]]{{end}}</title>
        {{block "styles" .}}{{end}}
        {{block "scripts" .}}{{end}}
    </head>
    <body>
        {{block "body" .}}
        {{end}}
    </body>
</html>