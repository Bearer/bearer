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
	"github.com/bearer/bearer/pkg/report/coverage"
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
		return &files.List{Coverage: incompleteCoverage()}, nil
	}

	var headFiles []files.File
	discovery := completeCoverage()

	gitFiles, err := git.ListTree(repository.context.RootDir, repository.context.CurrentCommitHash)
	if err != nil {
		return nil, err
	}

	for _, gitFile := range gitFiles {
		headFile, reason, fileErr := repository.fileForWithCoverage(
			ignore,
			goclocResult,
			gitFile.Filename,
		)
		if fileErr != nil {
			discovery.Errors = append(discovery.Errors, fmt.Sprintf("%s: %v", gitFile.Filename, fileErr))
			continue
		}
		if reason != "" {
			// Included paths and in-target exclusions share the fixed target namespace.
			relative, relErr := filepath.Rel(repository.targetPath, filepath.Join(repository.context.RootDir, gitFile.Filename))
			if relErr != nil {
				discovery.Errors = append(discovery.Errors, relErr.Error())
				continue
			}
			exclusion := coverage.Exclusion{Path: filepath.ToSlash(relative), Reason: reason}
			if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
				exclusion = coverage.Exclusion{Scope: "repository", Path: filepath.ToSlash(gitFile.Filename), Reason: "outside-target"}
			}
			discovery.Excluded = append(discovery.Excluded, exclusion)
			continue
		}
		if headFile != nil {
			headFiles = append(headFiles, *headFile)
			discovery.Files = append(discovery.Files, headFile.FilePath)
		}
	}

	discovery.Complete = len(discovery.Errors) == 0
	return &files.List{Files: headFiles, Coverage: discovery}, nil
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
		Coverage:  unsupportedCoverage(),
	}, nil
}

func unsupportedCoverage() coverage.Discovery {
	return coverage.Discovery{
		Files:    []string{},
		Excluded: []coverage.Exclusion{},
		Errors:   []string{"git diff discovery does not expose a supported same-run terminal file inventory"},
	}
}

func completeCoverage() coverage.Discovery {
	return coverage.Discovery{
		Files:    []string{},
		Excluded: []coverage.Exclusion{},
		Errors:   []string{},
		Complete: true,
	}
}

func incompleteCoverage() coverage.Discovery {
	coverage := completeCoverage()
	coverage.Complete = false
	coverage.Errors = append(coverage.Errors, "git current-tree discovery has no commit inventory")
	return coverage
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
	file, _, _ := repository.fileForWithCoverage(ignore, goclocResult, gitRelativePath)
	return file
}

func (repository *Repository) fileForWithCoverage(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
	gitRelativePath string,
) (*files.File, string, error) {
	relativePath, err := filepath.Rel(repository.gitTargetPath, gitRelativePath)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, "", err
	}

	if relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return nil, "outside-target", nil
	}

	fullPath := filepath.Join(repository.targetPath, relativePath)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		log.Debug().Msgf("error getting file stat: %s, %s", fullPath, err)
		return nil, "", err
	}

	if ignore != nil {
		reason, decisionErr := ignore.Decision(repository.targetPath, fullPath, goclocResult, fileInfo)
		if decisionErr != nil {
			return nil, "", decisionErr
		}
		if reason != "" {
			return nil, reason, nil
		}
	}

	return &files.File{
		Timeout:  timeout.Assign(fileInfo, repository.config),
		FilePath: relativePath,
	}, "", nil
}
