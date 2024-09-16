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
	action        string
	templateFiles []string
	outputFiles   []string
	enterprise    string
	githubToken   string
	validated     bool
	loaded        bool
	userList      UserList
	ownDomains    []string
}

type UserList struct {
	Updated    string     `json:"updated"`
	Enterprise Enterprise `json:"enterprise"`
	Users      []*User    `json:"users"`
	Warnings   []*Warning `json:"warnings"`
}

type Warning struct {
	Message string `json:"message"`
	Last    bool   `json:"last"`
}

type Enterprise struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type User struct {
	Number        int    `json:"number"`
	Login         string `json:"login"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	IsOwnDomain   bool   `json:"is_own_domain"`
	Contributions int    `json:"contributions"`
	Organizations *[]Organization
	Last          bool `json:"last"`
}

type Organization struct {
	Login        string        `json:"login"`
	Name         string        `json:"name"`
	Repositories *[]Repository `json:"repositories"`
	Last         bool          `json:"last"`
}

type Repository struct {
	Name string `json:"name"`
	Last bool   `json:"last"`
}

func (c *UserListConfig) Validate() error {
	if c.action == "" {
		return errors.New("Action is required")
	}
	if len(c.templateFiles) == 0 {
		return errors.New("Template is required")
	}
	if len(c.outputFiles) == 0 {
		return errors.New("Output File is required")
	}
	if c.enterprise == "" {
		return errors.New("Enterprise is required")
	}
	if c.githubToken == "" {
		return errors.New("Github Token is required")
	}
	if len(c.templateFiles) != len(c.outputFiles) {
		return fmt.Errorf("Template and Output Files must have the same length: %d != %d (%v, %v)", len(c.templateFiles), len(c.outputFiles), c.templateFiles, c.outputFiles)
	}

	c.validated = true
	slog.Debug("Validated userlist",
		"action", c.action,
		"enterprise", c.enterprise,
		"templateFiles", c.templateFiles,
		"githubToken", "***",
		"outputFiles", c.outputFiles,
		slog.Any("ownDomains", c.ownDomains))
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

	for i, templateFileName := range ul.templateFiles {
		outputFileName := ul.outputFiles[i]

		slog.Info("Rendering userlist", "templateFile", templateFileName, "outputFile", outputFileName)
		templateFile, err := os.ReadFile(templateFileName)
		if err != nil {
			slog.Error("Unable to read template file", "error", err, "file", templateFileName)
			return err
		}

		tmpl := template.Must(template.New("userlist").Parse(string(templateFile)))
		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, ul.userList)
		if err != nil {
			slog.Error("Unable to render userlist", "error", err)
			return err
		}

		err = os.WriteFile(outputFileName, buffer.Bytes(), 0644)
		if err != nil {
			slog.Error("Unable to write userlist", "error", err, "file", outputFileName)
			return err
		}
	}
	return nil
}

func (organization *Organization) RenderOutput(ctx context.Context, templateContent string) (string, error) {
	// render the organization to output
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
	// mark all eixsting users as last = false
	for i, _ := range ul.Users {
		ul.Users[i].Last = false
	}
	// mark the new user as last = true
	user.Last = true
	ul.Users = append(ul.Users, &user)
}

func (ul *UserList) findUser(login string) *User {
	for _, u := range ul.Users {
		if u.Login == login {
			return u
		}
	}
	return nil
}

func (ul *UserList) createUser(number int, login string, name string, email string, contributions int) *User {
	user := &User{
		Number:        number,
		Login:         login,
		Name:          name,
		Email:         email,
		Contributions: contributions,
		Organizations: new([]Organization),
	}
	ul.upsertUser(*user)
	return user
}

func (u *User) upsertOrganization(org Organization) {
	for _, o := range *u.Organizations {
		if o.Name == org.Name {
			// organization was found
			for _, repo := range *org.Repositories {
				o.upsertRepository(repo)
			}
			return
		}
	}
	slog.Debug("Upserting organization", "name", org.Name)
	// mark all existing organizations as last = false
	for i, _ := range *u.Organizations {
		(*u.Organizations)[i].Last = false
	}
	// mark the new organization as last = true
	org.Last = true
	*u.Organizations = append(*u.Organizations, org)
}

func (o *Organization) upsertRepository(repo Repository) {
	for _, r := range *o.Repositories {
		if r.Name == repo.Name {
			// repo was found
			return
		}
	}
	slog.Debug("Upserting repository", "name", repo.Name, "organization", o.Name)
	// mark all existing repositories as last = false
	for i := range *o.Repositories {
		(*o.Repositories)[i].Last = false
	}
	// mark the new repository as last = true
	repo.Last = true
	*o.Repositories = append(*o.Repositories, repo)
}

func (u *User) findOrganization(login string) *Organization {
	for _, o := range *u.Organizations {
		if o.Login == login {
			return &o
		}
	}
	return nil
}

func (u *User) createOrganization(login string, name string) *Organization {
	org := &Organization{
		Login:        login,
		Name:         name,
		Repositories: new([]Repository),
		Last:         false,
	}
	u.upsertOrganization(*org)
	return org
}

func (o *Organization) findRepository(name string) *Repository {
	for _, r := range *o.Repositories {
		if r.Name == name {
			return &r
		}
	}
	return nil
}

func (o *Organization) createRepository(name string) *Repository {
	repo := &Repository{
		Name: name,
	}
	o.upsertRepository(*repo)
	return repo
}

func (c *UserList) addWarning(warning string) {
	if c.Warnings == nil {
		c.Warnings = make([]*Warning, 0)
	}
	// mark all exisint warnings as last = false
	for _, w := range c.Warnings {
		w.Last = false
	}
	c.Warnings = append(c.Warnings, &Warning{Message: warning, Last: true})
}
