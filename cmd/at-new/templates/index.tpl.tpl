{{define "title"}}Welcome{{end}}
{{define "body"}}
<div style="display: flex;
            position: absolute;
            left: 0; right: 0;
            top: 0; bottom: 0;
            background-color: gray;
            justify-content: center;
            align-items: center;">

    <div style="flex: 0 0 auto;">
        <h1>Welcome to [[.Name]]</h1>
    </div>
</div>
{{end}}