package stabber

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/revlist"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"golang.org/x/mod/semver"
)

func GitLatestSemverTag(repo *git.Repository) (*plumbing.Reference, error) {
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	tag, err := tags.Next()
	if err != nil {
		return nil, err
	}

	tags.ForEach(func(r *plumbing.Reference) error {
		if semver.Compare(tag.Name().Short(), r.Name().Short()) == -1 {
			tag = r
		}
		return nil
	})

	return tag, nil
}

func GitCommitsBetween(repo *git.Repository, startHash, endHash plumbing.Hash) ([]*object.Commit, error) {
	ignoredHashes, err := revlist.Objects(repo.Storer, []plumbing.Hash{startHash}, []plumbing.Hash{})
	if err != nil {
		return nil, err
	}

	hashes, err := revlist.Objects(repo.Storer, []plumbing.Hash{endHash}, ignoredHashes)
	if err != nil {
		return nil, err
	}

	commits := make([]*object.Commit, 0)
	for _, hash := range hashes {
		commit, err := repo.CommitObject(hash)
		if errors.Is(err, plumbing.ErrObjectNotFound) {
			continue
		} else if err != nil {
			return nil, err
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

func GitHeadTag(repo *git.Repository) (*plumbing.Reference, error) {
	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}

	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	defer tags.Close()

	var tag *plumbing.Reference
	err = tags.ForEach(func(r *plumbing.Reference) error {
		revHash, err := repo.ResolveRevision(plumbing.Revision(r.Name()))
		if err != nil {
			return err
		}

		if *revHash != headRef.Hash() {
			return nil
		}

		tag = r
		return storer.ErrStop
	})

	return tag, err
}

func GitIsRefOn(repo *git.Repository, ref, target *plumbing.Reference) (bool, error) {
	// resolves both lightweight and annotated tags
	refHash, err := repo.ResolveRevision(plumbing.Revision(ref.Name()))
	if err != nil {
		return false, err
	}

	return *refHash == target.Hash(), nil
}
