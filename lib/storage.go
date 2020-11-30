package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type Inner struct {
	Counter int
	Content map[int]GithubData
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
		s.inner.Content = map[int]GithubData{}
	}

	return &s, nil
}

func (s *Storage) Delete(id int) {
	delete(s.inner.Content, id)
}

func (s *Storage) Update(id int, issue GithubData) {
	issue.Key  = id
	s.inner.Content[id] = issue
}

func (s *Storage) StoreData(issue GithubData) error {
	s.withId(func(id int) {
		issue.Key = id
		s.inner.Content[id] = issue
	})

	return nil
}

func (s *Storage) withId(f func(id int)) {
	nextId := s.inner.Counter + 1
	f(nextId)
	s.inner.Counter = nextId
}

func (s *Storage) LoadData() map[int]GithubData {
	return s.inner.Content
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

func (s Storage) Flush() {
	b, _ := json.MarshalIndent(s.inner, "", "  ")
	ioutil.WriteFile(s.path, b, os.ModePerm)
}
