package main

import (
	"flag"
	"github.com/prodyna/github-users/userlist"
	"log/slog"
	"os"
)

const (
	keyOrganization = "ORGANIZATION"
	keyGithubToken  = "GITHUB_TOKEN"
	keyTemplateFile = "TEMPLATE_FILE"
)

func main() {
	userlistCmd := flag.NewFlagSet("userlist", flag.ExitOnError)
	userlistEnterprise := userlistCmd.String("enterprise", "", "The GitHub Enterprise to query for repositories.")
	ueerlistGithubToken := userlistCmd.String("github-token", "", "The GitHub Token to use for authentication.")
	ueerlistTemplateFile := userlistCmd.String("template-file", "template/userlist.tpl", "The template file to use for rendering the result.")

	if len(os.Args) < 2 {
		slog.Error("expected command 'userlist'")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "userlist":
		userlistCmd.Parse(os.Args[2:])
		ul := userlist.New(
			userlist.WithEnterprise(*userlistEnterprise),
			userlist.WithGithubToken(*ueerlistGithubToken),
			userlist.WithTemplate(*ueerlistTemplateFile),
		)
		err := ul.Validate()
		if err != nil {
			slog.Error("Invalid configuration", "error", err)
			userlistCmd.PrintDefaults()
			os.Exit(1)
		}
		err = ul.Load()
		if err != nil {
			slog.Error("Unable to load userlist", "error", err)
			os.Exit(1)
		}
		err = ul.Render()
		if err != nil {
			slog.Error("Unable to render userlist", "error", err)
			os.Exit(1)
		}
	default:
		slog.Error("expected command")
		os.Exit(1)
	}

	//ctx := context.Background()
	//src := oauth2.StaticTokenSource(
	//	&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	//)
	//httpClient := oauth2.NewClient(ctx, src)
	//
	//client := githubv4.NewClient(httpClient)
	//
	///*
	//	{
	//	  enterprise(slug: "prodyna") {
	//	    ownerInfo {
	//	      samlIdentityProvider {
	//	        externalIdentities(after: null, first: 100) {
	//	          pageInfo {
	//	            hasNextPage
	//	            endCursor
	//	          }
	//	          edges {
	//	            node {
	//	              user {
	//	                login
	//	              }
	//	              samlIdentity {
	//	                nameId
	//	              }
	//	            }
	//	          }
	//	        }
	//	      }
	//	    }
	//	  }
	//	}
	//*/
	//var query struct {
	//	Enterprise struct {
	//		Slug      string
	//		Name      string
	//		OwnerInfo struct {
	//			SamlIdentityProvider struct {
	//				ExternalIdentities struct {
	//					PageInfo struct {
	//						HasNextPage bool
	//						EndCursor   githubv4.String
	//					}
	//					Edges []struct {
	//						Node struct {
	//							User struct {
	//								Login         string
	//								Organizations struct {
	//									Nodes []struct {
	//										Name string
	//									}
	//								}
	//							}
	//							SamlIdentity struct {
	//								NameId string
	//							}
	//						}
	//					}
	//				} `graphql:"externalIdentities(after: $after, first: $first)"`
	//			}
	//		}
	//		/*
	//			Members struct {
	//				TotalCount int
	//				Nodes      []struct {
	//					EnterpriseUserAccount struct {
	//						Login string
	//						Name  string
	//						User  struct {
	//							Login string
	//							Name  string
	//						}
	//					} `graphql:"... on EnterpriseUserAccount"`
	//					User struct {
	//						Login string
	//					} `graphql:"... on User"`
	//				}
	//				PageInfo struct {
	//					EndCursor   githubv4.String
	//					HasNextPage bool
	//				}
	//			} `graphql:"members(first: $first)"`
	//		*/
	//	} `graphql:"enterprise(slug: $slug)"`
	//}
	//
	//variables := map[string]interface{}{
	//	"slug":  githubv4.String("prodyna"),
	//	"first": githubv4.Int(3),
	//	"after": (*githubv4.String)(nil),
	//}
	//
	//err := client.Query(ctx, &query, variables)
	//if err != nil {
	//	slog.ErrorContext(ctx, "Unable to query", "error", err)
	//}
	//
	//// print json string of the query result
	//output, err := json.MarshalIndent(query, "", "  ")
	//if err != nil {
	//	slog.ErrorContext(ctx, "Unable to marshal query result", "error", err)
	//}
	//fmt.Print(string(output))
}
