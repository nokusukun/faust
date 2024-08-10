package docgen

const tmpl = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}} - API Documentation</title>
    <style>
        body { font-family: Arial, sans-serif; margin: auto; padding: 20px; max-width: 800px; }
        h1 { color: #333; }
        h2 { color: #555; }
        h3 { color: #777; }
        .endpoint { margin-bottom: 20px; margin-left: 1em }
        .method { font-weight: bold; color: #007BFF; }
        .parameters { margin-left: 20px; }
        .parameters li { margin-bottom: 5px; }
        .param-name { font-weight: bold; }
        .param-type { color: #555; font-style: italic; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <p>{{.Summary}}</p>
    <p><strong>Version:</strong> {{.Version}}</p>
    <h2>Routes</h2>

	{{range .Endpoints}}
	<div class="endpoint">
		<p class="method">{{.Method}} {{.Path}}</p>
		<p>{{.Description}}</p>
		{{if .Parameters}}
		<p><strong>Parameters:</strong></p>
		<ul class="parameters">
			{{range .Parameters}}
			<li><span class="param-name">{{.Name}}</span> (in {{.In}}) - {{.Description}} <span class="param-type">[{{.Schema.Type}}]</span></li>
			{{end}}
		</ul>
		{{end}}
	</div>
	{{end}}

    {{range .Subroutes}}
        <h3>Path: {{.Path}}</h3>
        {{range .Endpoints}}
        <div class="endpoint">
            <p class="method">{{.Method}} {{.Path}}</p>
            <p>{{.Description}}</p>
            {{if .Parameters}}
            <p><strong>Parameters:</strong></p>
            <ul class="parameters">
                {{range .Parameters}}
                <li><span class="param-name">{{.Name}}</span> (in {{.In}}) - {{.Description}} <span class="param-type">[{{.Schema.Type}}]</span></li>
                {{end}}
            </ul>
            {{end}}
        </div>
        {{end}}
    {{end}}
</body>
</html>
`
