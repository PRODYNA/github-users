# GitHub Enterprise Users for {{ .Enterprise.Name }}

Last updated: {{ .Updated }}

| # | GitHub Login | E-Mail |
| --- | --- | --- |
{{ range .Users }} | {{ .Number }} | [{{ .Login }}](https://github.com/enterprises/{{ $.Enterprise.Slug }}/people/{{ .Login }}/sso) | {{ .Email }} |
{{end}}
---
Generated with :heart: by [github-users](https://github.com/prodyna/github-users)
