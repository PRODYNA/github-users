# GitHub Enterprise Users

| # | GitHub Login | E-Mail |
| --- | --- | --- |
{{ range .Users }} | {{ .Number }} | [{{ .Login }}](https://github.com/enterprises/prodyna/people/{{ .Login }}/sso) | {{ .Email }} |
{{end}}
