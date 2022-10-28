package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/types"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wlredeye/jsonlines"
)

func ReportJSON(report types.Report, output *zerolog.Event) error {
	var ouputDetections any
	var err error

	if report.Type == flag.ReportDetectors {
		ouputDetections, err = getDetectorsOutput(report)
		if err != nil {
			return err
		}
	} else if report.Type == flag.ReportDataFlow {
		ouputDetections, err = getDataFlowOutput(report)
	}

	jsonBytes, err := json.MarshalIndent(&ouputDetections, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to json marshal detections: %w", err)
	}

	output.Msg(string(jsonBytes))

	return nil
}

func getDetectorsOutput(report types.Report) ([]interface{}, error) {
	var detections []interface{}
	f, err := os.Open(report.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open report: %w", err)
	}

	err = jsonlines.Decode(f, &detections)
	if err != nil {
		return nil, fmt.Errorf("failed to decode report: %w", err)
	}
	log.Debug().Msgf("got %d detections", len(detections))

	return detections, nil
}

// risks:
//   - detector_id: "rails_logger_detector"
//     data_types:
//       - name: "email"
//         locations:
//           - filename: "app/models/user.rb"
//             line_number: 5
//           - filename: "app/models/employee.rb"
//             line_number: 5
// risks:
//   - name: "email"
//	   detectors:
//		  - id: "rails_logger_detector"
//			stored: true
//          locations:
//				 - filename: "structure.sql"
//      	       line_number: 13
// data_types:
//   - name: "email"
//     detectors:
//       - name: sql_structure_detector
//         stored: true
//         locations:
//           - filename: "structure.sql"
//             line_number: 5
//             encrypted: true
//             verified_by:
//               - detector: "Rails Encryption"
//                 filename: "app/models/user.rb"
//                 line_number: 2

type DataFlow struct {
	Datatypes Datatype
}

type Datatype struct {
}

func getDataFlowOutput(report types.Report) (interface{}, error) {
	detections, err := getDetectorsOutput(report)
	if err != nil {
		return nil, err
	}

}
