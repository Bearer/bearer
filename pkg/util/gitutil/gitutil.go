package gitutil

import (
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/rs/zerolog/log"
)

func DiscoverFromGit(path string) (paths []string, err error) {
	// Instantiate a new repository targeting the .git directory
	fs := osfs.New(path)
	if _, err = fs.Stat(git.GitDirName); err == nil {
		fs, err = fs.Chroot(git.GitDirName)
		if err != nil {
			return
		}
	} else {
		log.Debug().Msg("No .git directory found")
		err = nil
		return
	}

	s := filesystem.NewStorageWithOptions(fs, cache.NewObjectLRUDefault(), filesystem.Options{KeepDescriptors: true})
	r, err := git.Open(s, fs)
	if err != nil {
		return
	}

	defer s.Close()

	ref, err := r.Head()
	if err != nil {
		return
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return
	}

	tree, err := commit.Tree()
	if err != nil {
		return
	}

	err = tree.Files().ForEach(func(f *object.File) error {
		// log.Debug().Msgf("100644 blob %s    %s", f.Hash, f.Name)
		paths = append(paths, f.Name)
		return nil
	})
	if err != nil {
		return
	}

	return
}
