package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(ctx, src)

	client := githubv4.NewClient(httpClient)

	/*
		{
		  enterprise(slug: "prodyna") {
		    ownerInfo {
		      samlIdentityProvider {
		        externalIdentities(after: null, first: 100) {
		          pageInfo {
		            hasNextPage
		            endCursor
		          }
		          edges {
		            node {
		              user {
		                login
		              }
		              samlIdentity {
		                nameId
		              }
		            }
		          }
		        }
		      }
		    }
		  }
		}
	*/
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
			/*
				Members struct {
					TotalCount int
					Nodes      []struct {
						EnterpriseUserAccount struct {
							Login string
							Name  string
							User  struct {
								Login string
								Name  string
							}
						} `graphql:"... on EnterpriseUserAccount"`
						User struct {
							Login string
						} `graphql:"... on User"`
					}
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"members(first: $first)"`
			*/
		} `graphql:"enterprise(slug: $slug)"`
	}

	variables := map[string]interface{}{
		"slug":  githubv4.String("prodyna"),
		"first": githubv4.Int(3),
		"after": (*githubv4.String)(nil),
	}

	err := client.Query(ctx, &query, variables)
	if err != nil {
		slog.ErrorContext(ctx, "Unable to query", "error", err)
	}

	// print json string of the query result
	output, err := json.MarshalIndent(query, "", "  ")
	if err != nil {
		slog.ErrorContext(ctx, "Unable to marshal query result", "error", err)
	}
	fmt.Print(string(output))
}
