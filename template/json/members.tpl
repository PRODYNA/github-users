{
    "enterprise": {
        "name": "{{ .Enterprise.Name }}",
        "slug": "{{ .Enterprise.Slug }}",
        "users": [{{ range .Users }}
            {
                "number": {{ .Number }},
                "login": "{{ .Login }}",
                "login_url": "https://github.com/enterprises/{{ $.Enterprise.Slug }}/people/{{ .Login }}/sso",
                "name": "{{ .Name }}",
                "email": "{{ .Email }}",
                "contributions": {{ .Contributions }},
                "is_own_domain": {{ .IsOwnDomain }}
            }{{ if not .Last }},{{ end }}{{ end }}
        ]
    },
    "warnings": [{{ range .Warnings }}
            "{{ . }}"{{ if not .Last }},{{ end }}
            {{ end }}
    ],
    "generated": {
        "at": "{{ .Updated }}",
        "by": "github-users",
        "with": ":heart:"
    }
}
