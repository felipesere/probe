package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

func (s *Storage) Update(id int, data GithubData) {
	data.Key = id
	s.inner.Content[id] = data
}

func (s *Storage) StoreData(data GithubData) {
	nextID := s.inner.Counter + 1
	data.Key = nextID
	s.inner.Content[nextID] = data

	s.inner.Counter = nextID
}

func (s *Storage) LoadData() map[int]GithubData {
	return s.inner.Content
}

func (s Storage) Flush() {
	b, _ := json.MarshalIndent(s.inner, "", "  ")
	_ = ioutil.WriteFile(s.path, b, os.ModePerm)
}
