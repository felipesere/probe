package gh

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/shurcooL/githubv4"
)

type Issue struct {
	status string
	id     string
}

func GetIssue(client githubv4.Client, owner, name string, nr int32) (Issue, error) {
	var query struct {
		Repository struct {
			Id    githubv4.ID
			Issue struct {
				State        githubv4.IssueState
				LastEditedAt githubv4.DateTime
				Comments     struct {
					Nodes []struct {
						BodyText githubv4.String
						Author   struct {
							Login string
						}
					}
				} `graphql:"comments(last: 1)"`
			} `graphql:"issue(number: $nr)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	err := client.Query(context.TODO(), &query,
		map[string]interface{}{
			"owner": githubv4.String(owner),
			"name":  githubv4.String(name),
			"nr":    githubv4.Int(nr),
		})

	issue := Issue{}
	if err != nil {
		return issue, err
	}

	spew.Dump(query)

	issue.status = fmt.Sprintf("%v", query.Repository.Issue.State)
	issue.id = fmt.Sprintf("%v", query.Repository.Id)

	return issue, nil
}
