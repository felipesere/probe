package gh

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"time"
)

type Issue struct {
	status      string
	id          string
	lastUpdated time.Time
	lastAction  string
}

type PR struct {
	status      string
	id          string
	lastUpdated time.Time
	lastAction  string
}

type CommentsQuery struct {
	Nodes []struct {
		BodyText githubv4.String
		Author   struct {
			Login string
		}
	}
}

type TimelineItemQuery struct {
	Nodes     []TimelineItem
	UpdatedAt time.Time
}

type TimelineItem struct {
	Typename string `graphql:"__typename"`
}

func GetPr(client githubv4.Client, owner, name string, nr int32) (PR, error) {
	var query struct {
		Repository struct {
			PullRequest struct {
				Id            githubv4.ID
				State         githubv4.PullRequestState
				LastEditedAt  githubv4.DateTime
				Comments      CommentsQuery     `graphql:"comments(last: 1)"`
				TimelineItems TimelineItemQuery `graphql:"timelineItems(last: 1)"`
			} `graphql:"pullRequest(number: $nr)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	err := client.Query(context.TODO(), &query,
		map[string]interface{}{
			"owner": githubv4.String(owner),
			"name":  githubv4.String(name),
			"nr":    githubv4.Int(nr),
		})

	pr := PR{}
	if err != nil {
		return pr, err
	}
	pr.status = fmt.Sprintf("%v", query.Repository.PullRequest.State)
	pr.id = fmt.Sprintf("%v", query.Repository.PullRequest.Id)
	pr.lastUpdated = query.Repository.PullRequest.TimelineItems.UpdatedAt
	pr.lastAction = query.Repository.PullRequest.TimelineItems.Nodes[0].Typename

	return pr, nil
}

func GetIssue(client githubv4.Client, owner, name string, nr int32) (Issue, error) {
	var query struct {
		Repository struct {
			Issue struct {
				Id            githubv4.ID
				State         githubv4.IssueState
				Comments      CommentsQuery `graphql:"comments(last: 1)"`
				TimelineItems TimelineItemQuery `graphql:"timelineItems(last: 1)"`
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

	issue.status = fmt.Sprintf("%v", query.Repository.Issue.State)
	issue.id = fmt.Sprintf("%v", query.Repository.Issue.Id)
	issue.lastUpdated = query.Repository.Issue.TimelineItems.UpdatedAt
	issue.lastAction = query.Repository.Issue.TimelineItems.Nodes[0].Typename

	return issue, nil
}
