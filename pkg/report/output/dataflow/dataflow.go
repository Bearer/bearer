package dataflow

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/customdetectors"
	"github.com/bearer/bearer/pkg/report/detections"
	reportdetectors "github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/output/dataflow/components"
	"github.com/bearer/bearer/pkg/report/output/dataflow/datatypes"
	"github.com/bearer/bearer/pkg/report/output/dataflow/detectiondecoder"
	fileerrors "github.com/bearer/bearer/pkg/report/output/dataflow/file_errors"
	"github.com/bearer/bearer/pkg/report/output/dataflow/paths"
	"github.com/bearer/bearer/pkg/report/output/dataflow/risks"
	"github.com/bearer/bearer/pkg/report/output/types"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/output"
)

var allowedDetections []detections.DetectionType = []detections.DetectionType{
	detections.TypeSchemaClassified,
	detections.TypeCustomClassified,
	detections.TypeDependencyClassified,
	detections.TypeInterfaceClassified,
	detections.TypeFrameworkClassified,
	detections.TypeCustomRisk,
	detections.TypeSecretleak,
	detections.TypeError,
	detections.TypeFileList,
	detections.TypeFileFailed,
	detections.TypeExpectedDetection,
	detections.TypeOperation,
}

func contains(detections []detections.DetectionType, detection detections.DetectionType) bool {
	for _, v := range detections {
		if v == detection {
			return true
		}
	}

	return false
}

func AddReportData(reportData *types.ReportData, config settings.Config, isInternal, hasFiles bool) error {
	if !hasFiles {
		reportData.Dataflow = &types.DataFlow{
			Languages: reportData.FoundLanguages,
		}
		return nil
	}

	expectedHolder := risks.New(config, isInternal)
	dataTypesHolder := datatypes.New(config, isInternal)
	risksHolder := risks.New(config, isInternal)
	componentsHolder := components.New(isInternal)
	pathsHolder := paths.New(isInternal)
	errorsHolder := fileerrors.New()

	extras, err := datatypes.NewExtras(reportData.Detectors, config)
	if err != nil {
		return err
	}

	customExtras, err := datatypes.NewCustomExtras(reportData.Detectors, config)
	if err != nil {
		return err
	}

	var files []string
	for _, detection := range reportData.Detectors {
		detectionMap, ok := detection.(map[string]interface{})
		if !ok {
			return fmt.Errorf("found detection in report which is not object")
		}

		detectionTypeS, ok := detectionMap["type"].(string)

		if !ok {
			continue
		}

		detectionType := detections.DetectionType(detectionTypeS)

		isDataflow := contains(allowedDetections, detectionType)
		if !isDataflow {
			continue
		}

		switch detectionType {
		case detections.TypeFileList:
			var fileListDetection detections.FileListDetection

			buf := bytes.NewBuffer(nil)
			if err := json.NewEncoder(buf).Encode(detection); err != nil {
				return err
			}
			if err = json.NewDecoder(buf).Decode(&fileListDetection); err != nil {
				return err
			}

			files = fileListDetection.Filenames
		case detections.TypeFileFailed:
			var errorDetection detections.FileFailedDetection
			buf := bytes.NewBuffer(nil)
			if err := json.NewEncoder(buf).Encode(detection); err != nil {
				return err
			}
			if err = json.NewDecoder(buf).Decode(&errorDetection); err != nil {
				return err
			}

			errorsHolder.AddFileError(errorDetection)
		case detections.TypeError:
			var errorDetection detections.ErrorDetection
			buf := bytes.NewBuffer(nil)
			if err := json.NewEncoder(buf).Encode(detection); err != nil {
				return err
			}
			if err = json.NewDecoder(buf).Decode(&errorDetection); err != nil {
				return err
			}

			errorsHolder.AddError(errorDetection)
		default:
			var castDetection detections.Detection

			buf := bytes.NewBuffer(nil)
			if err := json.NewEncoder(buf).Encode(detection); err != nil {
				return err
			}
			if err = json.NewDecoder(buf).Decode(&castDetection); err != nil {
				return err
			}

			// add full path to filename
			fullFilename := file.GetFullFilename(config.Target, castDetection.Source.Filename)
			castDetection.Source.FullFilename = fullFilename

			switch detectionType {
			case detections.TypeSchemaClassified:
				var detectionExtras *datatypes.ExtraFields
				if castDetection.DetectorType == reportdetectors.DetectorSchemaRb {
					detectionExtras = customExtras.Get(detection)
				}

				if err = dataTypesHolder.AddSchema(castDetection, detectionExtras); err != nil {
					return err
				}
			case detections.TypeOperation:
				operationDetection, err := detectiondecoder.GetOperation(detection)
				if err != nil {
					return err
				}
				pathsHolder.AddOperation(castDetection.DetectorType, operationDetection, fullFilename)
			case detections.TypeExpectedDetection:
				expectedHolder.AddRiskPresence(castDetection)
			case detections.TypeCustomRisk:
				ruleName := string(castDetection.DetectorType)
				customDetector, ok := config.Rules[ruleName]
				if !ok {
					customDetector, ok = config.BuiltInRules[ruleName]
					if !ok {
						return fmt.Errorf("custom detector not in config %s", ruleName)
					}
				}
				if customDetector.Type == customdetectors.TypeShared {
					continue
				}

				risksHolder.AddRiskPresence(castDetection)
			case detections.TypeCustomClassified:
				ruleName := string(castDetection.DetectorType)
				customDetector, ok := config.Rules[ruleName]
				if !ok {
					customDetector, ok = config.BuiltInRules[ruleName]
					if !ok {
						return fmt.Errorf("custom detector not in config %s", ruleName)
					}
				}

				switch customDetector.Type {
				case customdetectors.TypeVerifier:
				case customdetectors.TypeShared:
					continue
				case customdetectors.TypeRisk:
					if err := risksHolder.AddSchema(castDetection); err != nil {
						return err
					}
				case customdetectors.TypeDatatype:
					var detectionExtras *datatypes.ExtraFields
					detectionExtras = extras.Get(detection)

					if err = dataTypesHolder.AddSchema(castDetection, detectionExtras); err != nil {
						return err
					}
				}
			case detections.TypeSecretleak:
				risksHolder.AddRiskPresence(castDetection)
			case detections.TypeDependencyClassified:
				classifiedDetection, err := detectiondecoder.GetClassifiedDependency(detection)
				if err != nil {
					return err
				}

				classifiedDetection.Source.FullFilename = fullFilename
				err = componentsHolder.AddDependency(classifiedDetection)
				if err != nil {
					return err
				}
			case detections.TypeInterfaceClassified:
				classifiedDetection, err := detectiondecoder.GetClassifiedInterface(detection)
				if err != nil {
					return err
				}

				classifiedDetection.Source.FullFilename = fullFilename
				if err = componentsHolder.AddInterface(classifiedDetection); err != nil {
					return err
				}
			case detections.TypeFrameworkClassified:
				classifiedDetection, err := detectiondecoder.GetClassifiedFramework(detection)
				if err != nil {
					return err
				}
				classifiedDetection.Source.FullFilename = fullFilename
				if err = componentsHolder.AddFramework(classifiedDetection); err != nil {
					return err
				}
			}
		}
	}

	if !config.Scan.Quiet {
		output.StdErrLog("Generating dataflow")
	}

	reportData.Files = files
	reportData.Dataflow = &types.DataFlow{
		Languages:          reportData.FoundLanguages,
		Datatypes:          dataTypesHolder.ToDataFlow(),
		ExpectedDetections: expectedHolder.ToDataFlow(),
		Risks:              risksHolder.ToDataFlow(),
		Components:         componentsHolder.ToDataFlow(),
		Dependencies:       componentsHolder.ToDataFlowForDependencies(),
		Errors:             errorsHolder.ToDataFlow(),
		Paths:              pathsHolder.ToDataFlow(),
	}

	return nil
}
