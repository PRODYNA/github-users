# GitHub Enterprise collaborators for {{ .Enterprise.Name }}

Last updated: {{ .Updated }}

| Number | User | Contributions | Organization | Repository |
| ------ | ---- | ------------- | ------------ | ---------- |
{{ range $user := .Users }}{{ range $org := $user.Organizations }}{{ range $repo := $org.Repositories }}| {{ $user.Number }} | [{{ $user.Login }}](https://github.com/{{ $user.Login }}) | {{if $user.Contributions}}:green_square:{{else}}:red_square:{{end}} {{ $user.Contributions }} | [{{ $org.Name }}](https://github.com/{{ $org.Login }}) | [{{ $repo.Name }}](https://github.com/{{ $org.Login }}/{{ $repo.Name }}) |
{{ end }}{{ end }}{{ end }}

{{ if .Warnings }}
## Warnings
{{ range .Warnings }}* {{ . }}
{{ end }}{{ end }}
---
Generated with :heart: by [github-users](https://github.com/prodyna/github-users)
