package main

import (
	"text/template"
)

type tmplArguments struct {
	// Global options
	Options map[string]string
	Servers []string

	// Parameters related to current request
	AllServersLinkActive bool
	AllServersURL        string
	// Whois specific handling (for its unique URL)
	IsWhois     bool
	WhoisTarget string

	URLOption  string
	URLServer  string
	URLCommand string

	// Generated content to be displayed
	Title   string
	Content string
}

var tmpl = template.Must(template.New("tmpl").Parse(`
<!DOCTYPE html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
	<title>GCC Looking Glasses</title>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@4.5.0/dist/css/bootstrap.min.css" rel="stylesheet">
	<script src="https://cdn.jsdelivr.net/npm/jquery@3.2.1/dist/jquery.min.js"></script>
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@4.5.0/dist/js/bootstrap.min.js"></script>
	<script src="https://cdn.jsdelivr.net/npm/viz.js@2.1.2/viz.min.js" crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/viz.js@2.1.2/lite.render.js" crossorigin="anonymous"></script>
</head>
<body>

<div class="container">
	<nav class="navbar navbar-expand-lg navbar-light bg-light">
	    <a class="navbar-brand" href="/">GCC Bird Looking Glasses</a>
	    <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
	        <span class="navbar-toggler-icon"></span>
	    </button>

	    <div class="collapse navbar-collapse" id="navbarSupportedContent">
			<ul class="navbar-nav mr-auto">
				<li class="nav-item dropdown">
					<a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
				          {{ if .AllServersLinkActive }} All Servers {{ else }}{{ $.URLServer }}{{ end }}
			        </a>
					<div class="dropdown-menu" aria-labelledby="navbarDropdown">
						<a class="dropdown-item" href="/{{ .URLOption }}/{{ .AllServersURL }}/{{ .URLCommand }}"> All Servers </a>
						{{ range $k, $v := .Servers }}
						<a class="dropdown-item" href="/{{ $.URLOption }}/{{ $v }}/{{ $.URLCommand }}">{{ $v }}</a>
						{{ end }}
					</div>
				</li>
			</ul>
	    </div>

		{{ $option := .URLOption }}
		{{ $target := .URLCommand }}
		{{ if .IsWhois }}
			{{ $option = "whois" }}
			{{ $target = .WhoisTarget }}
		{{ end }}
		<form class="form-inline" action="/redir" method="GET">
			<div class="input-group">
				<select name="action" class="form-control">
					{{ range $k, $v := .Options }}
					<option value="{{ $k }}"{{ if eq $k $option }} selected{{end}}>{{ $v }}</option>
					{{ end }}
				</select>
				<input name="server" class="d-none" value="{{ .URLServer }}">
				<input name="target" class="form-control" placeholder="Target" aria-label="Target" value="{{ $target }}">
				<div class="input-group-append">
					<button class="btn btn-outline-success" type="submit">&raquo;</button>
				</div>
			</div>
		</form>

	</nav>

{{ .Content }}
</div>
</body>
</html>
`))
