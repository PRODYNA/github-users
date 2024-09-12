package userlist

import "strings"

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

func WithTemplateFile(templateFile string) func(*UserListConfig) {
	return func(config *UserListConfig) {
		config.templateFile = templateFile
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

func WithOutputFile(outputFile string) func(*UserListConfig) {
	return func(config *UserListConfig) {
		config.outputFile = outputFile
	}
}

func WithOwnDomains(ownDomains string) func(*UserListConfig) {
	return func(config *UserListConfig) {
		config.ownDomains = strings.Split(ownDomains, ",")
	}
}
