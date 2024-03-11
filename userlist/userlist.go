package userlist

import (
	"errors"
	"log/slog"
)

type UserListConfig struct {
	template    string
	enterprise  string
	githubToken string
	userList    *UserList
	validated   bool
	loaded      bool
}

type UserList struct {
	Users []User
}

type User struct {
	Number        int            `json:"Number"`
	Login         string         `json:"Login"`
	Email         string         `json:"Email"`
	Organizations []Organization `json:"Organizations"`
}

type Organization struct {
	Name string `json:"Name"`
}

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

func WithClient() func(*UserListConfig) {
	return func(config *UserListConfig) {
	}
}

func WithTemplate(template string) func(*UserListConfig) {
	return func(config *UserListConfig) {
		config.template = template
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

func (c *UserListConfig) Validate() error {
	if c.template == "" {
		return errors.New("Template is required")
	}
	if c.enterprise == "" {
		return errors.New("Enterprise is required")
	}
	if c.githubToken == "" {
		return errors.New("Github Token is required")
	}
	c.validated = true
	return nil
}

func (c *UserListConfig) Load() error {
	if !c.validated {
		return errors.New("Config not validated")
	}
	slog.Info("Loading userlist", "enterprise", c.enterprise)
	c.userList = &UserList{}
	c.loaded = true
	return nil
}

func (c *UserListConfig) Print() error {
	if !c.loaded {
		return errors.New("UserList not loaded")
	}
	return nil
}

func (ul *UserListConfig) Render() error {
	if !ul.loaded {
		return errors.New("UserList not loaded")
	}
	slog.Info("Rendering userlist", "template", ul.template)
	return nil
}
