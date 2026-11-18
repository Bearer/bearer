package gitleaks

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/secret"
	"github.com/bearer/bearer/pkg/report/source"
	"github.com/bearer/bearer/pkg/util/file"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
	"github.com/zricethezav/gitleaks/v8/sources"
)

//go:embed gitlab_config.toml
var RawConfig []byte

type detector struct {
	gitleaksDetector *detect.Detector
	config           config.Config
	idGenerator      nodeid.Generator
}

func New(idGenerator nodeid.Generator) types.Detector {
	var vc config.ViperConfig
	toml.Unmarshal(RawConfig, &vc) //nolint:all,errcheck
	cfg, err := vc.Translate()
	if err != nil {
		log.Fatal(err)
	}

	gitleaksDetector := detect.NewDetector(cfg)

	return &detector{
		gitleaksDetector: gitleaksDetector,
		config:           cfg,
		idGenerator:      idGenerator,
	}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	reader, err := file.Reader()
	if err != nil {
		return false, err
	}
	defer reader.Close()

	findings, err := detector.gitleaksDetector.DetectSource(
		context.Background(),
		&readerSource{
			content:         reader,
			Config:          &detector.config,
			Path:            file.AbsolutePath,
			MaxArchiveDepth: detector.gitleaksDetector.MaxArchiveDepth,
		},
	)
	if err != nil {
		return false, err
	}

	for _, finding := range findings {
		text := strings.TrimPrefix(finding.Line, "\n")
		report.AddSecretLeak(secret.Secret{
			Description: finding.Description,
		}, source.Source{
			Filename:          file.Path.RelativePath,
			StartLineNumber:   &finding.StartLine,
			StartColumnNumber: &finding.StartColumn,
			EndLineNumber:     &finding.EndLine,
			EndColumnNumber:   &finding.EndColumn,
			Text:              &text,
		})
	}

	return false, nil
}

// readerSource adapts an already-acquired reader to Gitleaks' single-file
// source contract. Path remains metadata used for archive suffix selection and
// outer allowlists; it is never opened by this source.
type readerSource struct {
	content         io.Reader
	Path            string
	Config          *config.Config
	MaxArchiveDepth int
}

var ErrArchiveExpansionUnsupported = errors.New("Bearer Gitleaks adapter: enabled archive expansion is unsupported")

func (s *readerSource) Fragments(ctx context.Context, yield sources.FragmentsFunc) error {
	// Reject the effective request before allowlists or content handling; silently
	// skipping an enabled expansion request would misrepresent this contract.
	if s.MaxArchiveDepth > 0 {
		return ErrArchiveExpansionUnsupported
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.Config != nil {
		for _, a := range s.Config.Allowlists {
			if a.PathAllowed(s.Path) {
				return nil
			}
		}
	}
	if s.content == nil {
		return fmt.Errorf("gitleaks input reader is nil")
	}

	input := &sources.File{
		Content:         s.content,
		Path:            s.Path,
		Config:          s.Config,
		MaxArchiveDepth: s.MaxArchiveDepth,
	}

	var fragmentErr error
	err := input.Fragments(ctx, func(fragment sources.Fragment, err error) error {
		if err != nil {
			// Preserve Gitleaks' normal behavior of scanning a fragment that
			// contains bytes, while making the source-level I/O failure visible
			// to DetectSource after its callback logs and continues.
			fragmentErr = err
		}
		if yieldErr := yield(fragment, err); yieldErr != nil {
			return yieldErr
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return fragmentErr
}
