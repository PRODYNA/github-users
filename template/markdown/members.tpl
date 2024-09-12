# GitHub Enterprise members for {{ .Enterprise.Name }}

Last updated: {{ .Updated }}

| # | GitHub Login | GitHub name | E-Mail | Contributions |
| --- | --- | --- | --- | --- |
{{ range .Users }} | {{ .Number }} | [{{ .Login }}](https://github.com/enterprises/{{ $.Enterprise.Slug }}/people/{{ .Login }}/sso) | {{ .Name }} | {{ if .IsOwnDomain }}:green_square:{{else}}:red_square:{{end}} {{ .Email }}  | {{if .Contributions}}:green_square:{{else}}:red_square:{{end}} [{{.Contributions }}](https://github.com/{{ .Login }}) |
{{ end }}

{{ if .Users }}_{{ len .Users }} users_{{ else }}No users found.{{ end }}

{{ if .Warnings }}
## Warnings
{{ range .Warnings }}* {{ . }}
{{ end }}{{ end }}
---
Generated with :heart: by [github-users](https://github.com/prodyna/github-users)
