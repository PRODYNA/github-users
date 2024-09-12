{
    "updated": "{{ .Updated }}",
    "enterprise": {
        "name": "{{ .Enterprise.Name }}",
        "slug": "{{ .Enterprise.Slug }}"
    },
    "users": [{{ range $user := .Users }}
        {
            "number": {{ $user.Number }},
            "login": "{{ $user.Login }}",
            "contributions": {{ $user.Contributions }},
            "organizations": [{{ range $org := $user.Organizations }}
                {
                    "name": "{{ $org.Name }}",
                    "login": "{{ $org.Login }}",
                    "repositories": [{{ range $repo := $org.Repositories }}
                        {
                            "name": "{{ $repo.Name }}"
                        }{{ if not $repo.Last }},{{ end }}{{ end }}
                    ]
                }{{ if not $org.Last }},{{ end }}{{ end }}
            ]
        }{{ if not $user.Last }},{{ end }}{{ end }}
    ],
    "warnings": [{{ range .Warnings }}
        "{{ .Message }}"{{ if not .Last }},{{ end }}{{ end }}
    ],
    "generated": {
        "by": "github-users",
        "with": ":heart:"
    }
}
