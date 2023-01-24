package gitleaks

import (
	_ "embed"
	"log"
	"regexp"

	"github.com/bearer/curio/pkg/detectors/types"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/secret"
	"github.com/bearer/curio/pkg/report/source"
	"github.com/bearer/curio/pkg/util/file"
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

// https://github.com/Skyscanner/whispers/blob/master/whispers/rules/sensitive-files.yml
var sensitiveFilesRegex = regexp.MustCompile(`.*(
	rsa|dsa|ed25519|ecdsa|pem|crt|cer|ca-bundle|p7b|p7c|p7s|(private-)?key|keystore|jks|pkcs12|pfx|p12|asc
	|dockercfg|npmrc|pypirc|pip.conf|terraform.tfvars|env|cfg|conf|config|ini|s3cfg
	|\.aws/credentials|htpasswd|(\.|-)?netrc|git-credentials|gitconfig|gitrobrc
	|(password|credential|secret)(\.[A-Za-z0-9]+)?
	|servlist-?\.conf|irssi/config|keys\.db
	|settings\.py|database\.yml
	|(config(\.inc)?|LocalSettings)\.php
	|(secret-token|omniauth|carrierwave|schema|knife)\.rb
	|(accounts|dbeaver-data-sources|BapSshPublisherPlugin|credentials|filezilla|recentservers)\.xml
	|(ba|z|da)?sh-history|(bash|zsh)rc|(bash-|zsh-)?(profile|aliases)
	|kdbx?|(agile)?keychain|key(store|ring)
	|mysql-history|psql-history|pgpass|irb-history
	|\.log|pcap|sql(dump)?|gnucash|dump
	|backup|back|bck|~1
	|kwallet|tblk|pubxml(\.user)?
	|Favorites\.plist|configuration\.user\.xpl
	|proftpdpasswd|robomongo\.json
	|ventrilo-srv.ini|muttrc|\.trc|ovpn|dayone|tugboat
  )$`)

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

	if sensitiveFilesRegex.Match([]byte(file.Name())) {
		lineNumber := 1
		columnNumber := 1
		report.AddSecretLeak(secret.Secret{
			Description: "Sensitive file name",
		}, source.Source{
			Filename:     file.Path.RelativePath,
			LineNumber:   &lineNumber,
			ColumnNumber: &columnNumber,
		})
		return false, nil
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

	if err != nil {
		return false, err
	}

	return false, nil
}
