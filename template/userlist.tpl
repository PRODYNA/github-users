# GitHub Enterprise Users for {{ .Enterprise.Name }}

| # | GitHub Login | E-Mail |
| --- | --- | --- |
{{ range .Users }} | {{ .Number }} | [{{ .Login }}](https://github.com/enterprises/{{ $.Enterprise.Slug }}/people/{{ .Login }}/sso) | {{ .Email }} |
{{end}}
