package main

import (
	"fmt"
	"log"

	"github.com/bearer/curio/new/builtin/detectors/ruby/objects"
	"github.com/bearer/curio/new/builtin/detectors/ruby/properties"
	"github.com/bearer/curio/new/detectionexecutor"
	"github.com/bearer/curio/new/detectioninitiator"
	"github.com/bearer/curio/new/language"
	"gopkg.in/yaml.v2"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	lang, err := language.Get("ruby")
	if err != nil {
		return fmt.Errorf("failed to lookup language: %s", err)
	}

	executor := detectionexecutor.New(lang)

	propertiesDetector, err := properties.New(lang)
	if err != nil {
		return fmt.Errorf("failed to create properties detector: %s", err)
	}
	err = executor.RegisterDetector(propertiesDetector)
	if err != nil {
		return fmt.Errorf("failed to register properties detector: %s", err)
	}

	objectsDetector, err := objects.New(lang)
	if err != nil {
		return fmt.Errorf("failed to create objects detector: %s", err)
	}
	err = executor.RegisterDetector(objectsDetector)
	if err != nil {
		return fmt.Errorf("failed to register objects detector: %s", err)
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

	initiator := detectioninitiator.New(executor, tree)

	detections, err := initiator.TreeDetections(tree.RootNode(), "objects")
	if err != nil {
		return fmt.Errorf("failed to detect objects: %s", err)
	}

	detectionsYAML, _ := yaml.Marshal(detections)
	log.Printf("GOT:\n%s\n", detectionsYAML)

	return nil
}
