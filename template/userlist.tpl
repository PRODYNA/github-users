# GitHub Enterprise Users for {{ .Enterprise.Name }}

Last updated: {{ .Updated }}

<<<<<<< Updated upstream
| # | GitHub Login | E-Mail |
| --- | --- | --- |
{{ range .Users }} | {{ .Number }} | [{{ .Login }}](https://github.com/enterprises/{{ $.Enterprise.Slug }}/people/{{ .Login }}/sso) | {{ .Email }} |
=======
| # | GitHub Login | GitHub name | E-Mail |
| --- | --- | --- | --- |
{{ range .Users }} | {{ .Number }} | [{{ .Login }}](https://github.com/enterprises/{{ $.Enterprise.Slug }}/people/{{ .Login }}/sso) | {{ .Name }} | {{ .Email }} |
>>>>>>> Stashed changes
{{end}}
---
Generated with :heart: by [github-users](https://github.com/prodyna/github-users)
