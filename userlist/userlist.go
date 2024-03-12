package userlist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"text/template"
)

const (
	members       = "members"
	collaborators = "collaborators"
)

type UserListConfig struct {
	action       string
	templateFile string
	markdownFile string
	enterprise   string
	githubToken  string
	validated    bool
	loaded       bool
	userList     UserList
}

type UserList struct {
	Updated    string
	Enterprise Enterprise
	Users      []*User
}

type Enterprise struct {
	Slug string
	Name string
}

type User struct {
	Number        int    `json:"Number"`
	Login         string `json:"Login"`
	Name          string `json:"Name"`
	Email         string `json:"Email"`
	Contributions int    `json:"Contributions"`
	Organizations *[]Organization
}

type Organization struct {
	Name         string        `json:"Name"`
	Repositories *[]Repository `json:"Repositories"`
}

type Repository struct {
	Name string `json:"Name"`
}

func (c *UserListConfig) Validate() error {
	if c.action == "" {
		return errors.New("Action is required")
	}
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
		"action", c.action,
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
	switch c.action {
	case members:
		return c.loadMembers()
	case collaborators:
		return c.loadCollaborators()
	default:
		return errors.New(fmt.Sprintf("Unknown action %s", c.action))
	}
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

func (ul *UserList) upsertUser(user User) {
	for i, u := range ul.Users {
		if u.Login == user.Login {
			*ul.Users[i] = user
			return
		}
	}
	slog.Info("Upserting user", "login", user.Login)
	ul.Users = append(ul.Users, &user)
}

func (u *User) upsertOrganization(org Organization) {
	//for i, o := range u.Organizations {
	//	if o.Name == org.Name {
	//		u.Organizations[i] = org
	//		return
	//	}
	//}
	*u.Organizations = append(*u.Organizations, org)
	slog.Debug("Upserting organization", "name", org.Name, "user", u.Login)
}

func (o *Organization) upsertRepository(repo Repository) {
	//for i, r := range o.Repositories {
	//	if r.Name == repo.Name {
	//		o.Repositories[i] = repo
	//		return
	//	}
	//}
	slog.Debug("Upserting repository", "name", repo.Name, "organization", o.Name)
	*o.Repositories = append(*o.Repositories, repo)
}
