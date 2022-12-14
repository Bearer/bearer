package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/bearer/curio/new/builtin/detectors/ruby/datatypes"
	"github.com/bearer/curio/new/builtin/detectors/ruby/objects"
	"github.com/bearer/curio/new/builtin/detectors/ruby/properties"
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
		user = { first_name: "hello", last_name: "there" }
		HTTP.post("http://api.com", user)

		user = { first_name: "hello2", last_name: "there2" }
		HTTP.post("http://api.com", user)

		HTTP.post("http://api.com", { x: user })

		x = { first_name: "hello3", last_name: "there3" }
		HTTP.post("http://api.com", { user: x })
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

	detectionsYAML, _ := yaml.Marshal(detections)
	log.Printf("GOT:\n%s\n", detectionsYAML)

	return nil
}
