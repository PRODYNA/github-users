package userlist

import "strings"

const (
	separator = ","
)

func New(options ...func(*UserListConfig)) *UserListConfig {
	config := &UserListConfig{
		validated: false,
		loaded:    false,
	}
	for _, option := range options {
		option(config)
	}
	return config
}

func WithAction(action string) func(*UserListConfig) {
	return func(config *UserListConfig) {
		config.action = action
	}
}

func WithTemplateFiles(templateFiles string) func(*UserListConfig) {
	return func(config *UserListConfig) {
		config.templateFiles = strings.Split(templateFiles, separator)
	}
}

func WithEnterprise(enterprise string) func(*UserListConfig) {
	return func(config *UserListConfig) {
		config.enterprise = enterprise
	}
}

func WithGithubToken(githubToken string) func(*UserListConfig) {
	return func(config *UserListConfig) {
		config.githubToken = githubToken
	}
}

func WithOutputFiles(outputFiles string) func(*UserListConfig) {
	return func(config *UserListConfig) {
		config.outputFiles = strings.Split(outputFiles, separator)
	}
}

func WithOwnDomains(ownDomains string) func(*UserListConfig) {
	return func(config *UserListConfig) {
		config.ownDomains = strings.Split(ownDomains, separator)
	}
}
