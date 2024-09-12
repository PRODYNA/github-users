package config

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"strconv"
)

const (
	keyAction                   = "action"
	kkeyActionEnvironment       = "ACTION"
	keyEnterprise               = "enterprise"
	keyEnterpriseEnvironment    = "ENTERPRISE"
	keyGithubToken              = "githubToken"
	keyGithubTokenEnvironment   = "GITHUB_TOKEN"
	keyTemplateFiles            = "template-files"
	keyTemplateFilesEnvironment = "TEMPLATE_FILES"
	keyOutputFiles              = "output-files"
	keyOutputFilesEnvironment   = "OUTPUT_FILES"
	keyVerbose                  = "verbose"
	keyVerboseEnvironment       = "VERBOSE"
	keyOwnDomains               = "own-domains"
	keyOwnDomainsEnvironment    = "OWN_DOMAINS"
)

type Config struct {
	Action        string
	Enterprise    string
	GithubToken   string
	TemplateFiles string
	OutputFiles   string
	OwnDomains    string
}

func New() (*Config, error) {
	c := Config{}
	flag.StringVar(&c.Action, keyAction, lookupEnvOrString(kkeyActionEnvironment, ""), "The action to perform.")
	flag.StringVar(&c.Enterprise, keyEnterprise, lookupEnvOrString(keyEnterpriseEnvironment, ""), "The GitHub Enterprise to query for repositories.")
	flag.StringVar(&c.GithubToken, keyGithubToken, lookupEnvOrString(keyGithubTokenEnvironment, ""), "The GitHub Token to use for authentication.")
	flag.StringVar(&c.TemplateFiles, keyTemplateFiles, lookupEnvOrString(keyTemplateFilesEnvironment, "template/members.tpl"), "The template file to use for rendering the result.")
	flag.StringVar(&c.OutputFiles, keyOutputFiles, lookupEnvOrString(keyOutputFilesEnvironment, ""), "The output file to write the result to.")
	flag.StringVar(&c.OwnDomains, keyOwnDomains, lookupEnvOrString(keyOwnDomainsEnvironment, ""), "The comma separated list of domains to consider as own domains.")
	verbose := flag.Int(keyVerbose, lookupEnvOrInt(keyVerboseEnvironment, 0), "Verbosity level, 0=info, 1=debug. Overrides the environment variable VERBOSE.")

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
