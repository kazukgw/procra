package procra

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
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

type HTMLS3Repository struct {
	BucketName string
	Region     aws.Region
	s3.ACL
	AccessKey string
	SecretKey string
	*aws.Auth
	*s3.Bucket
	client *s3.S3
}

func (repo *HTMLS3Repository) Init() error {
	auth, err := aws.GetAuth(repo.AccessKey, repo.SecretKey)
	if err != nil {
		return err
	}
	repo.Auth = auth
	repo.Client = s3.New(auth, repo.Region)
	repo.Bucket = repo.Client.Bucket(repo.BucketName)
	return nil
}

func (repo *HTMLS3Repository) Auth() {

}

func (repo *HTMLS3Repository) Put(key string, dat []byte) error {
	return repo.Bucket.Put(key, dat, "html", s3.Private)
}

func (repo *HTMLS3Repository) Get(key string) ([]byte, error) {
	return repo.Bucket.Get(key)
}
