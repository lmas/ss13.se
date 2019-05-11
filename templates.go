package ss13_se

import (
	"html/template"
)

func loadTemplates() (map[string]*template.Template, error) {
	tmpls := make(map[string]*template.Template)
	for name, src := range tmplList {
		t, err := parseTemplate(tmplBase, src)
		if err != nil {
			return nil, err
		}
		tmpls[name] = t
	}
	return tmpls, nil
}

func parseTemplate(src ...string) (*template.Template, error) {
	var err error
	t := template.New("*")
	for _, s := range src {
		t, err = t.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

// Using the awesome style from http://bettermotherfuckingwebsite.com/

const tmplBase string = `<!DOCTYPE html>
<html>
        <head>
                <meta charset="utf-8">
		<link rel="stylesheet" href="/static/style.css" type="text/css">
                <title>
                        {{block "title" .}}NO TITLE{{end}} | ss13.se
                </title>
        </head>
        <body>
                <header>
			<a href="/">ss13.se</a>
			<a href="/server/{{.Hub.ID}}">Global stats</a>
			<p class="right">Last updated: {{.Hub.LastUpdated}}</p>
                </header>

                <section id="body">
                        {{block "body" .}}NO BODY{{end}}
                </section>

                <footer>
			<a href="https://github.com/lmas/ss13_se">Source</a>
                </footer>
        </body>
</html>`

var tmplList = map[string]string{
	"style": `
* {
	padding: 0px;
	margin: 0px;
}
body {
	margin: 0px auto;
	max-width: 1024px;
	font-size: 18px;
	padding: 0 10px;
	line-height: 1.6;
	color: #444;
	background-color: #fff;
}
h1, h2 {
	text-align: center;
}
a, a:hover, a:visited {
	color: #444;
	text-decoration: none;
}
a:hover {
	color: #000;
}
img {
	display: block;
	margin: auto;
}
header {
	margin-bottom: 40px;
	padding: 10px 20px;
	color: #fff;
	background-color: #444;
	border-bottom-left-radius: 5px;
	border-bottom-right-radius: 5px;
}
header a, header a:hover, header a:visited {
	color: #fff;
	text-decoration: none;
	display: inline;
	padding-right: 40px;
}
footer {
	margin-top: 40px;
	padding: 10px;
	text-align: center;
	font-size: 12px;
}
.button a {
	background-color: #444;
	color: #fff;
	border-radius: 5px;
	padding: 5px 10px;
	text-decoration: none;
}
.button a:hover {
	background-color: #888;
}
.left {
	float: left;
}
.right {
	float: right;
}
.hide td, .hide a {
	color: #bbb;
}
`,
	"index": `{{define "title"}}Index{{end}}
{{define "body"}}
<h1>Servers</h1>
<table>
	<thead><tr>
		<td>Players</td>
		<td>Server</td>
	</tr></thead>

	<tbody>
	{{range .Servers}}
		<tr {{if lt .Players 1}}class="hide"{{end}}>
			<td>{{.Players}}</td>
			<td><a href="/server/{{.ID}}">{{.Title}}</a></td>
		</tr>
	{{else}}
		<tr><td>0</td><td>Sorry, no servers yet!</td></tr>
	{{end}}
	</tbody>
</table>
{{end}}
`,

	"server": `{{define "title"}}{{.Server.Title}}{{end}}
{{define "body"}}
<h1>{{.Server.Title}}</h1>

{{if .Server.SiteURL}}
	<span class="button"><a href="{{.Server.SiteURL}}">Website</a></span>
{{end}}

{{if .Server.ByondURL}}
	<span class="button"><a href="{{.Server.ByondURL}}">Join game</a></span>
{{end}}

<p>Current players: {{.Server.Players}}</p>

<h2>Daily History</h2>
<img src="/server/{{.Server.ID}}/daily" alt="Unable to show a pretty graph">
<h2>Weekly History</h2>
<img src="/server/{{.Server.ID}}/weekly" alt="Unable to show a pretty graph">
<h2>Average per day</h2>
<img src="/server/{{.Server.ID}}/averagedaily" alt="Unable to show a pretty graph">
<h2>Average per hour</h2>
<img src="/server/{{.Server.ID}}/averagehourly" alt="Unable to show a pretty graph">
{{end}}
`,
}
