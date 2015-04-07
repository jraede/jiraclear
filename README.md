Run this in a git project to delete branches that have a certain status in JIRA. The branch names need to contain the exact issue ID. E.g. (feature/ABCD-1234) would look check if issue ABCD-1234 is live in JIRA.

# Usage
`jiraclear --project {project} --url {url} --username {username} --password {password} --status {status}`

## Arguments
1. Project - the project prefix/key in JIRA (e.g. for tickets named ABCD-1234, this is "ABCD")
2. Url - The full URL to your JIRA installation, including protocol (http(s)), without the trailing slash
3. Username - Your JIRA username
4. Password - Your JIRA password
5. Status - The status to check. Defaults to "Live"