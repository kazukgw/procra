package procra

import (
	"io/ioutil"
	"path"
)

type HTMLRepository interface {
	Get(key string) ([]byte, error)
	Put(key string, dat []byte) error
}

type HTMLFileRepository struct {
	Base string
}

func (repo HTMLFileRepository) Put(key string, dat []byte) error {
	return ioutil.WriteFile(path.Join(repo.Base, key, ".html"), dat, 0644)
}

func (repo HTMLFileRepository) Get(key string) ([]byte, error) {
	return ioutil.ReadFile(path.Join(repo.Base, key, ".html"))
}
