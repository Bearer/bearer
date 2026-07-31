package gitleaks

import (
	"context"
	_ "embed"
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
	findings, err := detector.gitleaksDetector.DetectSource(
		context.Background(),
		&sources.Files{
			Config:          &detector.config,
			Path:            file.AbsolutePath,
			Sema:            detector.gitleaksDetector.Sema,
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
