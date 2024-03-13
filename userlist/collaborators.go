package userlist

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"log/slog"
	"time"
)

func (c *UserListConfig) loadCollaborators() error {
	slog.Info("Loading collaborators", "enterprise", c.enterprise)
	c.userList = UserList{
		// updated as RFC3339 string
		Updated: time.Now().Format(time.RFC3339),
	}
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.githubToken},
	)
	httpClient := oauth2.NewClient(ctx, src)
	client := githubv4.NewClient(httpClient)

	/*
		{
		  enterprise(slug: "prodyna") {
		    slug
		    name
		    organizations (first:100) {
		      nodes {
				login
		        name
		      }
		    }
		  }
		}
	*/
	var organizations struct {
		Enterprise struct {
			Slug          string
			Name          string
			Organizations struct {
				Nodes []struct {
					Login string
					Name  string
				}
			} `graphql:"organizations(first:100)"`
		} `graphql:"enterprise(slug: $slug)"`
	}

	variables := map[string]interface{}{
		"slug": githubv4.String(c.enterprise),
	}

	slog.Info("Loading organizations", "enterprise", c.enterprise)
	err := client.Query(ctx, &organizations, variables)
	if err != nil {
		slog.ErrorContext(ctx, "Unable to query", "error", err)
		return err
	}
	slog.Info("Loaded organizations", "organization.count", len(organizations.Enterprise.Organizations.Nodes))

	/*
		{
		  organization(login:"prodyna") {
		    repositories(first:100) {
		      pageInfo {
		        hasNextPage
		        startCursor
		      }
		      nodes {
		        name
		        collaborators(first:100,affiliation:OUTSIDE) {
		          pageInfo {
		            hasNextPage
		            startCursor
		          }
		          nodes {
		            login
		            name
		          }
		        }
		      }
		    }
		  }
		}
	*/
	c.userList.Enterprise.Slug = organizations.Enterprise.Slug
	c.userList.Enterprise.Name = organizations.Enterprise.Name

	userNumber := 0
	slog.Info("Iterating organizatons", "organization.count", len(organizations.Enterprise.Organizations.Nodes))
	for _, org := range organizations.Enterprise.Organizations.Nodes {
		if org.Login != "PRODYNA" {
			continue
		}
		slog.Info("Loading repositories and external collaborators", "organization", org.Login)
		var query struct {
			Organization struct {
				Login        string
				Repositories struct {
					Nodes []struct {
						Name          string
						Collaborators struct {
							Nodes []struct {
								Login                   string
								Name                    string
								ContributionsCollection struct {
									ContributionCalendar struct {
										TotalContributions int
									}
								}
							}
						} `graphql:"collaborators(first:100,affiliation:OUTSIDE)"`
					}
				} `graphql:"repositories(first:100)"`
			} `graphql:"organization(login: $organization)"`
		}

		variables := map[string]interface{}{
			"organization": githubv4.String(org.Login),
		}

		err := client.Query(ctx, &query, variables)
		if err != nil {
			slog.WarnContext(ctx, "Unable to query - will skip this organization", "error", err, "organization", org.Login)
			continue
		}

		// count the collaborators
		collaboratorCount := 0
		for _, repo := range query.Organization.Repositories.Nodes {
			collaboratorCount += len(repo.Collaborators.Nodes)
		}
		if collaboratorCount == 0 {
			slog.DebugContext(ctx, "No collaborators found", "organization", org.Login)
			continue
		}

		for _, repo := range query.Organization.Repositories.Nodes {
			slog.DebugContext(ctx, "Processing repository", "repository", repo.Name, "collaborator.count", len(repo.Collaborators.Nodes))
			for _, collaborator := range repo.Collaborators.Nodes {
				slog.DebugContext(ctx, "Processing collaborator", "login", collaborator.Login, "name", collaborator.Name, "contributions", collaborator.ContributionsCollection.ContributionCalendar.TotalContributions)
				user := c.userList.findUser(collaborator.Login)
				if user == nil {
					user = c.userList.createUser(userNumber+1, collaborator.Login, collaborator.Name, "", collaborator.ContributionsCollection.ContributionCalendar.TotalContributions)
					userNumber++
				} else {
					slog.Info("Found existing user", "login", user.Login)
				}
				organization := Organization{
					Login:        org.Login,
					Name:         org.Name,
					Repositories: new([]Repository),
				}
				user.upsertOrganization(organization)
				repository := Repository{
					Name: repo.Name,
				}
				organization.upsertRepository(repository)
			}
		}

		slog.InfoContext(ctx, "Loaded repositories",
			"repository.count", len(query.Organization.Repositories.Nodes),
			"organization", org.Login,
			"collaborator.count", collaboratorCount)

		if collaboratorCount == 0 {
			continue
		}

		output, err := json.MarshalIndent(c.userList, "", "  ")
		if err != nil {
			slog.ErrorContext(ctx, "Unable to marshal json", "error", err)
			return err
		}
		fmt.Printf("%s\n", output)

		slog.InfoContext(ctx, "Adding collaborators", "organization", org.Login)
	}

	c.loaded = true
	return nil
}
