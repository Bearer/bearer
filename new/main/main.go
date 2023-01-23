package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/bearer/curio/new/builtin/detectors/ruby/datatypes"
	"github.com/bearer/curio/new/builtin/detectors/ruby/objects"
	"github.com/bearer/curio/new/builtin/detectors/ruby/properties"
	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/detectorexecutor"
	getlanguage "github.com/bearer/curio/new/language/get"
	"github.com/bearer/curio/new/treeevaluator"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	lang, err := getlanguage.Get("ruby")
	if err != nil {
		return fmt.Errorf("failed to lookup language: %s", err)
	}

	propertiesDetector, err := properties.New(lang)
	if err != nil {
		return fmt.Errorf("failed to create properties detector: %s", err)
	}
	defer propertiesDetector.Close()

	objectsDetector, err := objects.New(lang)
	if err != nil {
		return fmt.Errorf("failed to create objects detector: %s", err)
	}
	defer objectsDetector.Close()

	datatypesDetector, err := datatypes.New(lang)
	if err != nil {
		return fmt.Errorf("failed to create datatypes detector: %s", err)
	}
	defer datatypesDetector.Close()

	executor, err := detectorexecutor.New(lang, []detector.Detector{
		propertiesDetector,
		objectsDetector,
		datatypesDetector,
	})
	if err != nil {
		return fmt.Errorf("failed to create executor: %s", err)
	}

	tree, err := lang.Parse(`

class Abc
	def f
		user = { first_name: "hello", last_name: "there" }   # 5
		HTTP.post("http://api.com", user)                    # 6

		user = { first_name: "hello2", last_name: "there2" } # 8
		HTTP.post("http://api.com", user)                    # 9

		HTTP.post("http://api.com", { x: user })             # 11

		x = { first_name: "hello3", last_name: "there3" }    # 13
		HTTP.post("http://api.com", { user: x })             # 14
	end
end

	`)
	if err != nil {
		return err
	}
	defer tree.Close()

	evaluator := treeevaluator.New(lang, executor, tree)

	detections, err := evaluator.TreeDetections(tree.RootNode(), "datatypes")
	if err != nil {
		return fmt.Errorf("failed to detect data types: %s", err)
	}

	detectionsYAML, _ := yaml.Marshal(report(detections))
	log.Printf("GOT:\n%s\n", detectionsYAML)

	return nil
}

type Location struct {
	Content    string
	LineNumber int
}
type ReportDetection struct {
	MatchLocation   Location
	ContextLocation *Location
	Data            interface{}
}

func report(detections []*detectiontypes.Detection) []ReportDetection {
	result := make([]ReportDetection, len(detections))

	for i, detection := range detections {
		var contextLocation *Location

		if detection.ContextNode != nil {
			contextLocation = &Location{
				Content:    detection.ContextNode.Content(),
				LineNumber: detection.ContextNode.LineNumber(),
			}
		}

		result[i] = ReportDetection{
			MatchLocation: Location{
				Content:    detection.MatchNode.Content(),
				LineNumber: detection.MatchNode.LineNumber(),
			},
			Data:            detection.Data,
			ContextLocation: contextLocation,
		}
	}

	return result
}
