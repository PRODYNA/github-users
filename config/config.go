package config

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"strconv"
)

const (
	keyAction       = "ACTION"
	keyEnterprise   = "ENTERPRISE"
	keyGithubToken  = "GITHUB_TOKEN"
	keyTemplateFile = "TEMPLATE_FILE"
	keyVerbose      = "VERBOSE"
	keyMarkdownFile = "MARKDOWN_FILE"
)

type Config struct {
	Action       string
	Enterprise   string
	GithubToken  string
	TemplateFile string
	MarkdownFile string
}

func New() (*Config, error) {
	c := Config{}
	flag.StringVar(&c.Action, keyAction, lookupEnvOrString("ACTION", ""), "The action to perform.")
	flag.StringVar(&c.Enterprise, keyEnterprise, lookupEnvOrString("ENTERPRISE", ""), "The GitHub Enterprise to query for repositories.")
	flag.StringVar(&c.GithubToken, keyGithubToken, lookupEnvOrString("GITHUB_TOKEN", ""), "The GitHub Token to use for authentication.")
	flag.StringVar(&c.TemplateFile, keyTemplateFile, lookupEnvOrString("TEMPLATE_FILE", "template/members.tpl"), "The template file to use for rendering the result.")
	flag.StringVar(&c.MarkdownFile, keyMarkdownFile, lookupEnvOrString("MARKDOWN_FILE", "USERS.md"), "The markdown file to write the result to.")
	verbose := flag.Int("verbose", lookupEnvOrInt(keyVerbose, 0), "Verbosity level, 0=info, 1=debug. Overrides the environment variable VERBOSE.")

	level := slog.LevelInfo
	if *verbose > 0 {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})))
	flag.Parse()
	return &c, nil
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func lookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}
