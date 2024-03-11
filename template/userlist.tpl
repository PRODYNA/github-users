# GitHub Enterprise Users for {{ .Enterprise.Name }}

Last updated: {{ .Updated }}

| # | GitHub Login | GitHub name | E-Mail |
| --- | --- | --- |
{{ range .Users }} | {{ .Number }} | [{{ .Login }}](https://github.com/enterprises/{{ $.Enterprise.Slug }}/people/{{ .Login }}/sso) | {{ .Name }} | {{ .Email }} |
{{end}}
---
Generated with :heart: by [github-users](https://github.com/prodyna/github-users)
