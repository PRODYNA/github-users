package userlist

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

func WithMarkdownFile(markdownFile string) func(*UserListConfig) {
	return func(config *UserListConfig) {
		config.markdownFile = markdownFile
	}
}
