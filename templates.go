package ss13hub

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

const tmplBase string = `<!DOCTYPE html>
<html>
        <head>
                <meta charset="utf-8">
                <title>
                        {{block "title" .}}NO TITLE{{end}} | ss13.se
                </title>
                <style type="text/css">
                img {
			display: block;
			margin: auto;
		}
                footer {
			text-align: center;
		}
		.left {
			float: left;
		}
		.right {
			float: right;
		}
                </style>
        </head>
        <body>
                <header>
			<h2><a href="/">SS13.se</a></h2>
                </header>

                <section id="body">
                        {{block "body" .}}NO BODY{{end}}
                </section>

                <footer>
			<p>
				<span class="left">
					Source code at
					<a href="https://github.com/lmas/ss13_se">Github</a>
				</span>

				{{/* TODO: not sure about the copyright stuff when fetching ext. data */}}
				Copyright © 2017 A. Svensson

				<span class="right">
					Using raw data from
					<a href="http://www.byond.com/games/exadv1/spacestation13">Byond</a>
				</span>
			</p>
                </footer>
        </body>
</html>`

var tmplList = map[string]string{
	"index": `{{define "title"}}Index{{end}}
{{define "body"}}
<p>Last updated: {{.Hub.LastUpdated}}</p>
<p>Current # of servers: {{.TotalServers}}</p>
<p>Current # of players: {{.Hub.Players}}</p>
<a href="/server/{{.Hub.ID}}">Global stats</a><br />
<br />
<table>
	<thead><tr>
		<td>Players</td>
		<td>Server</td>
	</tr></thead>

	<tbody>
	{{range .Servers}}
		<tr>
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

<p>Last updated: {{.Server.LastUpdated}}</p>
<p>Current players: {{.Server.Players}}</p>
{{if .Server.SiteURL}}
	<a href="{{.Server.SiteURL}}">Web site</a><br />
{{end}}

{{if .Server.ByondURL}}
	<a href="{{.Server.ByondURL}}">Join game</a><br />
{{end}}

<br />
<img src="/server/{{.Server.ID}}/daily" alt="Daily history"><br />
<img src="/server/{{.Server.ID}}/weekly" alt="Weekly history"><br />
<img src="/server/{{.Server.ID}}/average" alt="Average per day"><br />
{{end}}
`,
}
