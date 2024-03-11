package userlist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"log/slog"
	"os"
	"text/template"
)

type UserListConfig struct {
	templateFile string
	markdownFile string
	enterprise   string
	githubToken  string
	userList     *UserList
	validated    bool
	loaded       bool
}

type UserList struct {
	Enterprise Enterprise
	Users      []User
}

type Enterprise struct {
	Slug string
	Name string
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

func (c *UserListConfig) Validate() error {
	if c.templateFile == "" {
		return errors.New("Template is required")
	}
	if c.enterprise == "" {
		return errors.New("Enterprise is required")
	}
	if c.githubToken == "" {
		return errors.New("Github Token is required")
	}
	if c.markdownFile == "" {
		return errors.New("Markdown File is required")
	}
	c.validated = true
	slog.Debug("Validated userlist",
		"enterprise", c.enterprise,
		"template", c.templateFile,
		"githubToken", "***",
		"markdownFile", c.markdownFile)
	return nil
}

func (c *UserListConfig) Load() error {
	if !c.validated {
		return errors.New("Config not validated")
	}
	slog.Info("Loading userlist", "enterprise", c.enterprise)
	c.userList = &UserList{}

	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.githubToken},
	)
	httpClient := oauth2.NewClient(ctx, src)
	client := githubv4.NewClient(httpClient)

	var query struct {
		Enterprise struct {
			Slug      string
			Name      string
			OwnerInfo struct {
				SamlIdentityProvider struct {
					ExternalIdentities struct {
						PageInfo struct {
							HasNextPage bool
							EndCursor   githubv4.String
						}
						Edges []struct {
							Node struct {
								User struct {
									Login         string
									Organizations struct {
										Nodes []struct {
											Name string
										}
									} `graphql:"organizations(first: 10)"`
								}
								SamlIdentity struct {
									NameId string
								}
							}
						}
					} `graphql:"externalIdentities(after: $after, first: $first)"`
				}
			}
		} `graphql:"enterprise(slug: $slug)"`
	}

	window := 100
	variables := map[string]interface{}{
		"slug":  githubv4.String("prodyna"),
		"first": githubv4.Int(window),
		"after": (*githubv4.String)(nil),
	}

	for offset := 0; ; offset += window {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			slog.ErrorContext(ctx, "Unable to query", "error", err)
		}

		c.userList.Enterprise = Enterprise{
			Slug: query.Enterprise.Slug,
			Name: query.Enterprise.Name,
		}

		for i, e := range query.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.Edges {
			u := User{
				Number: offset + i + 1,
				Login:  e.Node.User.Login,
				Email:  e.Node.SamlIdentity.NameId,
			}
			for _, o := range e.Node.User.Organizations.Nodes {
				u.Organizations = append(u.Organizations, Organization{Name: o.Name})
			}
			c.userList.Users = append(c.userList.Users, u)
		}

		if !query.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.PageInfo.HasNextPage {
			break
		}

		variables["after"] = githubv4.NewString(query.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.PageInfo.EndCursor)
	}

	slog.InfoContext(ctx, "Loaded userlist", "users", len(c.userList.Users))
	c.loaded = true
	return nil
}

func (c *UserListConfig) Print() error {
	if !c.loaded {
		return errors.New("UserList not loaded")
	}
	slog.Info("Printing userlist")
	output, err := json.MarshalIndent(c.userList, "", "  ")
	if err != nil {
		slog.Error("Unable to marshal json", "error", err)
		return err
	}
	fmt.Printf("%s\n", output)

	return nil
}

func (ul *UserListConfig) Render() error {
	if !ul.loaded {
		return errors.New("UserList not loaded")
	}
	slog.Info("Rendering userlist", "template", ul.templateFile)
	templateFile, err := os.ReadFile(ul.templateFile)
	if err != nil {
		slog.Error("Unable to read template file", "error", err, "file", ul.templateFile)
		return err
	}

	tmpl := template.Must(template.New("userlist").Parse(string(templateFile)))
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, ul.userList)
	if err != nil {
		slog.Error("Unable to render userlist", "error", err)
		return err
	}

	err = os.WriteFile(ul.markdownFile, buffer.Bytes(), 0644)
	if err != nil {
		slog.Error("Unable to write userlist", "error", err, "file", ul.markdownFile)
		return err
	}
	return nil
}

func (organization *Organization) RenderMarkdown(ctx context.Context, templateContent string) (string, error) {
	// render the organization to markdown
	tmpl := template.Must(template.New("organization").Parse(templateContent))
	// execute template to a string
	var buffer bytes.Buffer
	err := tmpl.Execute(&buffer, organization)
	if err != nil {
		slog.ErrorContext(ctx, "Unable to render organization", "error", err)
		return "", err
	}
	return buffer.String(), nil
}
