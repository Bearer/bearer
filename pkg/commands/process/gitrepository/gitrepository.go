package gitrepository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/filelist/ignore"
	"github.com/bearer/bearer/pkg/commands/process/filelist/timeout"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/git"
	"github.com/bearer/bearer/pkg/util/output"
)

type Repository struct {
	ctx    context.Context
	config settings.Config
	targetPath,
	gitTargetPath string
	context *Context
}

func New(ctx context.Context, config settings.Config, targetPath string, context *Context) (*Repository, error) {
	if context == nil {
		log.Debug().Msg("no git repository found")
		return nil, nil
	}

	gitTargetPath, err := filepath.Rel(context.RootDir, targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative target: %w", err)
	}

	log.Debug().Msgf("git target: [%s]/%s", context.RootDir, gitTargetPath)

	repository := &Repository{
		ctx:           ctx,
		config:        config,
		targetPath:    targetPath,
		gitTargetPath: gitTargetPath,
		context:       context,
	}

	if err = repository.fetchMergeBaseCommit(); err != nil {
		return nil, err
	}

	return repository, nil
}

func (repository *Repository) ListFiles(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
) (*files.List, error) {
	if repository == nil {
		return nil, nil
	}

	if repository.context.BaseCommitHash == "" {
		return repository.getCurrentFiles(ignore, goclocResult)
	}

	return repository.getDiffFiles(ignore, goclocResult)
}

func (repository *Repository) fetchMergeBaseCommit() error {
	hash := repository.context.BaseCommitHash
	if hash == "" {
		return nil
	}

	log.Debug().Msgf("merge base commit: %s", hash)

	if isPresent, err := git.CommitPresent(repository.context.RootDir, hash); isPresent || err != nil {
		return err
	}

	log.Debug().Msgf("merge base commit not present, fetching")

	if err := git.FetchRef(repository.ctx, repository.context.RootDir, hash); err != nil {
		return err
	}

	log.Debug().Msgf("merge base commit fetched")

	return nil
}

func (repository *Repository) getCurrentFiles(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
) (*files.List, error) {
	if repository.context.CurrentCommitHash == "" {
		return &files.List{}, nil
	}

	var headFiles []files.File

	gitFiles, err := git.ListTree(repository.context.RootDir, repository.context.CurrentCommitHash)
	if err != nil {
		return nil, err
	}

	for _, file := range gitFiles {
		if file := repository.fileFor(
			ignore,
			goclocResult,
			file.Filename,
		); file != nil {
			headFiles = append(headFiles, *file)
		}
	}

	return &files.List{Files: headFiles}, nil
}

func (repository *Repository) getDiffFiles(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
) (*files.List, error) {
	var baseFiles, headFiles []files.File
	renames := make(map[string]string)
	chunks := make(map[string]git.Chunks)

	filePatches, err := git.Diff(repository.context.RootDir, repository.context.BaseCommitHash)
	if err != nil {
		return nil, err
	}

	for _, patch := range filePatches {
		// we're not interested in removals
		if patch.ToPath == "" {
			continue
		}

		headFile := repository.fileFor(ignore, goclocResult, patch.ToPath)
		if headFile == nil {
			continue
		}

		headFiles = append(headFiles, *headFile)

		if patch.FromPath == "" {
			continue
		}

		relativeFromPath, err := filepath.Rel(repository.gitTargetPath, patch.FromPath)
		if err != nil {
			return nil, err
		}
		baseFiles = append(baseFiles, files.File{
			Timeout:  headFile.Timeout,
			FilePath: relativeFromPath,
		})

		if relativeFromPath != headFile.FilePath {
			renames[relativeFromPath] = headFile.FilePath
		}

		chunks[headFile.FilePath] = patch.Chunks
	}

	return &files.List{
		Files:     headFiles,
		BaseFiles: baseFiles,
		Renames:   renames,
		Chunks:    chunks,
	}, nil
}

func (repository *Repository) WithBaseBranch(body func() error) error {
	if repository == nil || !repository.config.Scan.Diff {
		return nil
	}

	if repository.context.HasUncommittedChanges {
		return errors.New("uncommitted changes found in your repository. commit or stash changes your changes and retry")
	}

	if err := git.Switch(repository.context.RootDir, repository.context.BaseCommitHash, true); err != nil {
		return fmt.Errorf("error switching to base branch: %w", err)
	}

	err := body()

	if restoreErr := repository.restoreCurrent(); restoreErr != nil {
		wrappedErr := fmt.Errorf("error restoring to current commit: %w", restoreErr)
		if err == nil {
			return wrappedErr
		}

		output.StdErrLog(wrappedErr.Error())
	}

	return err
}

func (repository *Repository) restoreCurrent() error {
	if repository.context.CurrentBranch == "" {
		return git.Switch(repository.context.RootDir, repository.context.CurrentCommitHash, true)
	}

	return git.Switch(repository.context.RootDir, repository.context.CurrentBranch, false)
}

func (repository *Repository) fileFor(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
	gitRelativePath string,
) *files.File {
	relativePath, err := filepath.Rel(repository.gitTargetPath, gitRelativePath)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	if strings.Contains(relativePath, "..") {
		return nil
	}

	fullPath := filepath.Join(repository.targetPath, relativePath)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		log.Debug().Msgf("error getting file stat: %s, %s", fullPath, err)
		return nil
	}

	if ignore != nil && ignore.Ignore(repository.targetPath, fullPath, goclocResult, fileInfo) {
		return nil
	}

	return &files.File{
		Timeout:  timeout.Assign(fileInfo, repository.config),
		FilePath: relativePath,
	}
}
