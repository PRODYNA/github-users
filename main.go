package main

import (
	config "github.com/prodyna/github-users/config"
	"github.com/prodyna/github-users/userlist"
	"log/slog"
	"os"
)

func main() {
	c, err := config.New()
	if err != nil {
		slog.Error("Unable to create config", "error", err)
		os.Exit(1)
	}

	ulc := userlist.New(
		userlist.WithAction(c.Action),
		userlist.WithEnterprise(c.Enterprise),
		userlist.WithGithubToken(c.GithubToken),
		userlist.WithTemplateFile(c.TemplateFile),
		userlist.WithOutputFile(c.OutputFile),
		userlist.WithOwnDomains(c.OwnDomains),
	)

	err = ulc.Validate()
	if err != nil {
		slog.Error("Invalid config", "error", err)
		os.Exit(1)
	}
	err = ulc.Load()
	if err != nil {
		slog.Error("Unable to load userlist", "error", err)
		os.Exit(1)
	}
	err = ulc.Print()
	if err != nil {
		slog.Error("Unable to print userlist", "error", err)
		os.Exit(1)
	}
	err = ulc.Render()
	if err != nil {
		slog.Error("Unable to render userlist", "error", err)
		os.Exit(1)
	}
}
