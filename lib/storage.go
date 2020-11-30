package lib

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type Row struct {
	GithubID    string
	Owner       string
	Repository  string
	Title       string
	Status      string
	LastAction  string
	LastChanged string
	Link        string
	Kind        string
	Number      int32
}

type Inner struct {
	Counter int
	Content map[int]Row
}

type Storage struct {
	path  string
	inner Inner
}

func NewStorage(path string) (*Storage, error) {
	b, _ := ioutil.ReadFile(path)
	s := Storage{path: path}
	err := json.Unmarshal(b, &s.inner)
	if err != nil {
		s.inner.Content = map[int]Row{}
	}

	return &s, nil
}

func (s *Storage) Delete(id int) {
	delete(s.inner.Content, id)
}

func (s *Storage) Update(id int, issue GithubData)  {
	s.inner.Content[id] = newItem(issue)
}

func newItem(issue GithubData) Row {
	return Row {
		GithubID:    issue.Id,
		Owner:       issue.Owner,
		Repository:  issue.Repository,
		Number:      issue.Number,
		Title:       issue.Title,
		Status:      issue.Status,
		LastAction:  issue.LastAction,
		LastChanged: issue.LastUpdated.Format(time.RFC3339),
		Link:        issue.Link,
		Kind:        string(issue.Kind),
	}
}

func (s *Storage) StoreData(issue GithubData) error {
	s.withId(func(id int) {
		s.inner.Content[id] = newItem(issue)
	})

	return nil
}

func (s *Storage) withId(f func(id int)) {
	nextId := s.inner.Counter + 1
	f(nextId)
	s.inner.Counter = nextId
}

func (s *Storage) LoadData() map[int]GithubData {
	data := make(map[int]GithubData, 0)
	for id, row := range s.inner.Content {
		 data[id] = GithubData{
			Kind:        GithubKind(row.Kind),
			Owner:       row.Owner,
			Repository:  row.Repository,
			Number:      row.Number,
			Link:        row.Link,
			Title:       row.Title,
			Status:      row.Status,
			Id:          row.GithubID,
			LastUpdated: mustParse(time.RFC3339, row.LastChanged),
			LastAction:  row.LastAction,
		}
	}

	return data
}

func mustParse(layout, t string) time.Time {
	res, err := time.Parse(layout, t)
	if err != nil {
		panic(err)
	}

	return res
}

func (s *Storage) innerLoad() {
	file, _ := os.OpenFile(s.path, os.O_RDONLY, os.ModePerm)

	err := json.NewDecoder(file).Decode(&s.inner)
	if err != nil {
		panic(err.Error())
	}
}

func (s *Storage) Load() [][]string {
	s.innerLoad() // Is this OK? Keep on reloading? Well, we are not interactive... so shrug?

	var keys []int
	for key, _ := range s.inner.Content {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	var data [][]string
	for _, k := range keys {
		row := s.inner.Content[k]
		t, _ := time.Parse(time.RFC3339, row.LastChanged)
		data = append(data, []string{
			fmt.Sprintf("%d", k),
			row.Owner,
			row.Repository,
			row.Title,
			row.Status,
			row.LastAction,
			humanize.Time(t),
			row.Link,
		})
	}

	return data
}

func (s Storage) Flush() {
	b, _ := json.MarshalIndent(s.inner, "", "  ")
	ioutil.WriteFile(s.path, b, os.ModePerm)
}
