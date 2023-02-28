package gitleaks

import (
	_ "embed"
	"log"

	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/secret"
	"github.com/bearer/bearer/pkg/report/source"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/pelletier/go-toml"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
)

//go:embed gitlab_config.toml
var RawConfig []byte

type detector struct {
	gitleaksDetector *detect.Detector
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
		idGenerator:      idGenerator,
	}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	findings, err := detector.gitleaksDetector.DetectFiles(file.Path.AbsolutePath)

	if err != nil {
		return false, err
	}

	for _, finding := range findings {
		report.AddSecretLeak(secret.Secret{
			Description: finding.Description,
		}, source.Source{
			Filename:     file.Path.RelativePath,
			LineNumber:   &finding.StartLine,
			ColumnNumber: &finding.StartColumn,
		})
	}

	return false, nil
}
