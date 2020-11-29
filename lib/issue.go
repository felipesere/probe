package lib

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"time"
)

type Issue struct {
	Title       string
	Status      string
	Id          string
	LastUpdated time.Time
	LastAction  string
}

type PR struct {
	Title       string
	Status      string
	Id          string
	LastUpdated time.Time
	LastAction  string
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
				Title         string
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
	pullRequest := query.Repository.PullRequest

	pr.Status = fmt.Sprintf("%v", pullRequest.State)
	pr.Id = fmt.Sprintf("%v", pullRequest.Id)
	pr.LastUpdated = pullRequest.TimelineItems.UpdatedAt
	pr.LastAction = pullRequest.TimelineItems.Nodes[0].Typename
	pr.Title = pullRequest.Title

	return pr, nil
}

func GetIssue(client githubv4.Client, owner, name string, nr int32) (Issue, error) {
	var query struct {
		Repository struct {
			Issue struct {
				Id            githubv4.ID
				Title         string
				State         githubv4.IssueState
				Comments      CommentsQuery     `graphql:"comments(last: 1)"`
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

	i := query.Repository.Issue
	issue.Status = fmt.Sprintf("%v", i.State)
	issue.Id = fmt.Sprintf("%v", i.Id)
	issue.LastUpdated = i.TimelineItems.UpdatedAt
	issue.LastAction = i.TimelineItems.Nodes[0].Typename
	issue.Title = i.Title

	return issue, nil
}
