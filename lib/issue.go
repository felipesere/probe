package lib

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"time"
)

type GithubKind string

const (
	PullRequestKind GithubKind = "PR"
	IssueKind       GithubKind = "Issue"
)

type GithubData struct {
	Key         int // should this be a *int to show that it might not be there?
	Kind        GithubKind
	Owner       string
	Repository  string
	Number      int32
	Link        string
	Title       string
	Status      string
	Id          string
	LastUpdated time.Time
	LastAction  string
}

type ItemQuery struct {
	Id            githubv4.ID
	Title         string
	State         githubv4.PullRequestState
	LastEditedAt  githubv4.DateTime
	Comments      CommentsQuery     `graphql:"comments(last: 1)"`
	TimelineItems TimelineItemQuery `graphql:"timelineItems(last: 1)"`
	Url     string
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

// find a better name
type Target struct {
	Owner string
	Name  string
	Nr    int32
}

func GetPr(client githubv4.Client, target Target) (GithubData, error) {
	var query struct {
		Repository struct {
			PullRequest ItemQuery `graphql:"pullRequest(number: $nr)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	err := client.Query(context.TODO(), &query,
		map[string]interface{}{
			"owner": githubv4.String(target.Owner),
			"name":  githubv4.String(target.Name),
			"nr":    githubv4.Int(target.Nr),
		})

	if err != nil {
		return GithubData{}, err
	}
	pr := from(query.Repository.PullRequest)
	pr.Kind = PullRequestKind
	pr.Owner = target.Owner
	pr.Repository = target.Name
	pr.Number = target.Nr

	return pr, nil
}

func GetIssue(client githubv4.Client, target Target) (GithubData, error) {
	var query struct {
		Repository struct {
			Issue ItemQuery `graphql:"issue(number: $nr)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	err := client.Query(context.TODO(), &query,
		map[string]interface{}{
			"owner": githubv4.String(target.Owner),
			"name":  githubv4.String(target.Name),
			"nr":    githubv4.Int(target.Nr),
		})

	if err != nil {
		return GithubData{}, err
	}

	i := query.Repository.Issue

	issue := from(i)
	issue.Kind = IssueKind
	issue.Owner = target.Owner
	issue.Repository = target.Name
	issue.Number = target.Nr

	return issue, nil
}

func from(i ItemQuery) GithubData {
	return GithubData{
		Status:      fmt.Sprintf("%v", i.State),
		Id:          fmt.Sprintf("%v", i.Id),
		LastUpdated: i.TimelineItems.UpdatedAt,
		LastAction:  i.TimelineItems.Nodes[0].Typename,
		Title:       i.Title,
		Link:        i.Url,
	}
}
