package userlist

import (
	"context"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"log/slog"
	"strings"
	"time"
)

func (c *UserListConfig) loadMembers() error {
	slog.Info("Loading members", "enterprise", c.enterprise)
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
									Login                   string
									Name                    string
									ContributionsCollection struct {
										ContributionCalendar struct {
											TotalContributions int
										}
									}
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

	window := 25
	variables := map[string]interface{}{
		"slug":  githubv4.String(c.enterprise),
		"first": githubv4.Int(window),
		"after": (*githubv4.String)(nil),
	}

	for offset := 0; ; offset += window {
		slog.Debug("Running query", "offset", offset, "window", window)
		err := client.Query(ctx, &query, variables)
		if err != nil {
			slog.ErrorContext(ctx, "Unable to query", "error", err)
			return err
		}

		c.userList.Enterprise = Enterprise{
			Slug: query.Enterprise.Slug,
			Name: query.Enterprise.Name,
		}

		for i, e := range query.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.Edges {
			u := User{
				Number:        offset + i + 1,
				Login:         e.Node.User.Login,
				Name:          e.Node.User.Name,
				Email:         e.Node.SamlIdentity.NameId,
				IsOwnDomain:   IsOwnDomain(e.Node.SamlIdentity.NameId, c.ownDomains),
				Contributions: e.Node.User.ContributionsCollection.ContributionCalendar.TotalContributions,
			}
			c.userList.upsertUser(u)
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

func IsOwnDomain(email string, ownDomains []string) bool {
	if len(ownDomains) == 0 {
		return true
	}
	for _, domain := range ownDomains {
		if strings.HasSuffix(email, domain) {
			return true
		}
	}
	return false
}
