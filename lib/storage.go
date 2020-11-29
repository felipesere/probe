package lib

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"os"
	"time"
)

type Row struct {
	Owner       string
	Repository  string
	Title       string
	Status      string
	LastAction  string
	LastChanged string
	Link        string
}

type Storage struct {
	path    string
	content map[string]Row
}

func NewStorage(path string) (Storage, error) {
	b, _ := ioutil.ReadFile(path)
	s := Storage{path: path}
	err := json.Unmarshal(b, &s.content)
	if err != nil {
		s.content = map[string]Row{}
	}

	return s, nil
}

func (s Storage) AddIssue(issue Issue) error {
	s.content[issue.Id] = Row{
		Owner:       issue.Owner,
		Repository:  issue.Repository,
		Title:       issue.Title,
		Status:      issue.Status,
		LastAction:  issue.LastAction,
		LastChanged: issue.LastUpdated.Format(time.RFC3339),
		Link:        "",
	}

	return nil
}

func (s Storage) AddPullRequest(pullRequest PR) error {
	s.content[pullRequest.Id] = Row{
		Owner:       pullRequest.Owner,
		Repository:  pullRequest.Repository,
		Title:       pullRequest.Title,
		Status:      pullRequest.Status,
		LastAction:  pullRequest.LastAction,
		LastChanged: pullRequest.LastUpdated.Format(time.RFC3339),
		Link:        "",
	}

	return nil
}

func (s Storage) Load() [][]string  {
	file, _ := os.OpenFile(s.path, os.O_RDONLY, os.ModePerm)

	err := json.NewDecoder(file).Decode(&s.content)
	if err != nil {
		panic(err.Error())
	}
	var data [][]string
	var idx int
	for _,  row := range s.content {
		t, _ := time.Parse(time.RFC3339,row.LastChanged)
		data = append(data, []string{
			fmt.Sprintf("%d", idx),
			row.Owner,
			row.Repository,
			row.Title,
			row.Status,
			row.LastAction,
			humanize.Time(t),
			row.Link,
		})
		idx++;
	}

	return data
}

func (s Storage) Flush()  {
	b, _ := json.Marshal(s.content)
	ioutil.WriteFile(s.path, b, os.ModePerm)
}
