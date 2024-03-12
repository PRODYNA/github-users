package userlist

import (
	"context"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"log/slog"
	"time"
)

func (c *UserListConfig) loadCollaborators() error {
	slog.Info("Loading collaborators", "enterprise", c.enterprise)
	c.userList = &UserList{
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
	slog.Info("Iterating organizatons", "organization.count", len(organizations.Enterprise.Organizations.Nodes))
	for _, org := range organizations.Enterprise.Organizations.Nodes {
		if org.Login != "PRODYNA" {
			continue
		}
		slog.Info("Loading repositories and external collaborators", "organization", org.Login)
		var repositories struct {
			Organization struct {
				Repositories struct {
					Nodes []struct {
						Name          string
						Collaborators struct {
							Nodes []struct {
								Login string
								Name  string
							}
						} `graphql:"collaborators(first:100,affiliation:OUTSIDE)"`
					}
				} `graphql:"repositories(first:100)"`
			} `graphql:"organization(login: $organization)"`
		}

		variables := map[string]interface{}{
			"organization": githubv4.String(org.Login),
		}

		err := client.Query(ctx, &repositories, variables)
		if err != nil {
			slog.WarnContext(ctx, "Unable to query - will skip this organization", "error", err, "organization", org.Login)
			continue
		}

		// count the collaborators
		collaboratorCount := 0
		for _, repo := range repositories.Organization.Repositories.Nodes {
			collaboratorCount += len(repo.Collaborators.Nodes)
		}

		slog.InfoContext(ctx, "Loaded repositories",
			"repository.count", len(repositories.Organization.Repositories.Nodes),
			"organization", org.Login,
			"collaborator.count", collaboratorCount)

		if collaboratorCount == 0 {
			continue
		}

		slog.InfoContext(ctx, "Adding collaborators", "organization", org.Login)
	}

	return nil
}
