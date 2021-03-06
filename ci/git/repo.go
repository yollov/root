package git

import (
	"errors"
	"os"

	gogit "gopkg.in/src-d/go-git.v4"
)

type Repository struct {
	Path   string
	Origin string
	gogit.Repository
}

func GetRepo(origin, path string) (*Repository, error) {
	cloneOpts := &gogit.CloneOptions{
		URL:           origin,
		ReferenceName: "refs/heads/master",
		SingleBranch:  true,
		Depth:         2, // TODO dynamic
	}

	if path == "" {
		return nil, errors.New("Wrong cloning path")
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// exists -> delete folder
		if err := os.RemoveAll(path); err != nil {
			return nil, err
		}
	}
	// not exists -> create dir
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	// clone last n commits
	repo, err := gogit.PlainClone(path, false, cloneOpts)
	if err != nil {
		return nil, err
	}

	return &Repository{path, origin, *repo}, nil
}

func (repo *Repository) Destroy() error {
	return os.RemoveAll(repo.Path)
}

func (repo *Repository) Diff() ([]string, error) {
	// get current commit
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	current, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	// get diff
	stats, err := current.Stats()
	if err != nil {
		return nil, err
	}

	paths := make([]string, len(stats))
	for i, fstat := range stats {
		paths[i] = fstat.Name
	}

	return paths, nil
}
