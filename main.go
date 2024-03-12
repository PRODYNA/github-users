package main

import (
	config "github.com/prodyna/github-users/config"
	"github.com/prodyna/github-users/userlist"
	"log/slog"
	"os"
)

const (
	keyAction       = "ACTION"
	keyOrganization = "ORGANIZATION"
	keyGithubToken  = "GITHUB_TOKEN"
	keyTemplateFile = "TEMPLATE_FILE"
)

type Config struct {
	Action       string
	Enterprise   string
	GithubToken  string
	TemplateFile string
}

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
		userlist.WithMarkdownFile(c.MarkdownFile),
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
