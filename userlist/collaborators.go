package userlist

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"log/slog"
	"time"
)

const windowSize = 100

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
				PageInfo struct {
					HasNextPage bool
					EndCursor   githubv4.String
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

	if organizations.Enterprise.Organizations.PageInfo.HasNextPage {
		slog.Warn("More organizations available - not yet implemented")
		c.userList.addWarning("More organizations available - not yet implemented")
	}

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
							PageInfo struct {
								HasNextPage bool
								EndCursor   githubv4.String
							}
						} `graphql:"collaborators(first:100,affiliation:OUTSIDE)"`
					}
					PageInfo struct {
						HasNextPage bool
						EndCursor   githubv4.String
					}
				} `graphql:"repositories(first:$first,after:$after)"`
			} `graphql:"organization(login: $organization)"`
		}

		variables := map[string]interface{}{
			"organization": githubv4.String(org.Login),
			"first":        githubv4.Int(20),
			"after":        (*githubv4.String)(nil),
		}

		for {
			err := client.Query(ctx, &query, variables)
			if err != nil {
				slog.WarnContext(ctx, "Unable to query - will skip this organization", "error", err, "organization", org.Login)
				c.userList.addWarning(fmt.Sprintf("Unable to query organization %s", org.Login))
				break
			}

			for _, repo := range query.Organization.Repositories.Nodes {
				slog.DebugContext(ctx, "Processing repository", "repository", repo.Name, "collaborator.count", len(repo.Collaborators.Nodes))
				for _, collaborator := range repo.Collaborators.Nodes {
					slog.DebugContext(ctx, "Processing collaborator", "login", collaborator.Login, "name", collaborator.Name, "contributions", collaborator.ContributionsCollection.ContributionCalendar.TotalContributions)

					// User
					user := c.userList.findUser(collaborator.Login)
					if user == nil {
						user = c.userList.createUser(userNumber+1, collaborator.Login, collaborator.Name, "", collaborator.ContributionsCollection.ContributionCalendar.TotalContributions)
						userNumber++
					} else {
						slog.Info("Found existing user", "login", user.Login)
					}

					// Organization
					organization := user.findOrganization(org.Login)
					if organization == nil {
						organization = user.createOrganization(org.Login, org.Name)
					} else {
						slog.Info("Found existing organization", "organization", organization.Name)
					}

					// Repository
					repository := organization.findRepository(repo.Name)
					if repository == nil {
						repository = organization.createRepository(repo.Name)
					} else {
						slog.Info("Found existing repository", "repository", repository.Name)
					}
					organization.upsertRepository(*repository)
				}
			}

			slog.InfoContext(ctx, "Loaded repositories",
				"repository.count", len(query.Organization.Repositories.Nodes),
				"organization", org.Login)

			if !query.Organization.Repositories.PageInfo.HasNextPage {
				break
			}

			slog.Info("More repositories available", "organization", org.Login, "after", query.Organization.Repositories.PageInfo.EndCursor)
			variables["after"] = githubv4.NewString(query.Organization.Repositories.PageInfo.EndCursor)
		}
	}

	c.loaded = true
	return nil
}
