# GitHub Enterprise collaborators for {{ .Enterprise.Name }}

Last updated: {{ .Updated }}

| Number | User | Contributions | Organization | Repository |
| ------ | ---- | ------------- | ------------ | ---------- |
{{ range $user := .Users }}{{ range $org := $user.Organizations }}{{ range $repo := $org.Repositories }}| {{ $user.Number }} | {{ $user.Login }} | {{ $user.Contributions }} | {{ $org.Name }} | {{ $repo.Name }} |
{{ end }}{{ end }}{{ end }}

---
Generated with :heart: by [github-users](https://github.com/prodyna/github-users)
