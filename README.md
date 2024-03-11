# GitHub Users

GitHub Action that can create various list of users from GitHub

* All known GitHub users with their linked SSO accounts
* List of all external collaborators and who invited them

## User list

Creates a markdown file with the list of all users of the enterprise.

### Example

> # GitHub Enterprise Users
> | # | GitHub Login                                                 | E-Mail          |
> | --- |--------------------------------------------------------------|-----------------|
> | 1 | [foo](https://github.com/enterprises/octocat/people/foo/sso) | foo@octocat.com |
> | 2 | [bar](https://github.com/enterprises/octocat/people/bar/sso) | bar@octocat.com |

### Using

This action can be used in a workflow like this:

```yaml
name: Create userlist

on:
  workflow_dispatch:
  # Every day at 07:00
  schedule:
    - cron: '0 7 * * *'

jobs:
  create-userlist:
    name: "Create userlist"
    runs-on: ubuntu-latest
    steps:
      # Checkout the existing content of thre repository
      - name: Checkout
        uses: actions/checkout@v2

      # Run the deployment overview action
      - name: Github users
        uses: prodyna/github-users@v0.3
        with:
          # The action to run
          action: userlist
          # The GitHub Enterprise to query for repositories
          enterprise: octocat
          # The GitHub Token to use for authentication
          github-token: ${{ secrets.GITHUB_TOKEN }}
          # The template file to use for rendering the result
          template-file: template/userlist.tpl
          # The markdown file to write the result to
          markdown-file: USERS.md
          # Verbosity level, 0=info, 1=debug
          verbose: 1

      # Push the generated files
      - name: Commit changes
        run: |
          git config --local user.email "darko@krizic.net"
          git config --local user.name "Deployment Overview"
          git add profile
          git commit -m "Add/update deployment overview"
```
