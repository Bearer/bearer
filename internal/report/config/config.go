package config

import (
	"encoding/json"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Regexp struct {
	regexp.Regexp
}

func RegexpMustCompile(value string) *Regexp {
	return &Regexp{Regexp: *regexp.MustCompile(value)}
}

func (r *Regexp) UnmarshalJSON(data []byte) error {
	var valueString string
	if err := json.Unmarshal(data, &valueString); err != nil {
		return fmt.Errorf("expected string regex: %s", err)
	}

	regexpValue, err := regexp.Compile(valueString)
	if err != nil {
		return fmt.Errorf("invalid regex: %s", err)
	}

	*r = Regexp{*regexpValue}

	return nil
}

func (r *Regexp) UnmarshalYAML(value *yaml.Node) error {
	var valueString string
	if err := value.Decode(&valueString); err != nil {
		return fmt.Errorf("expected string regex: %s", err)
	}

	regexpValue, err := regexp.Compile(valueString)
	if err != nil {
		return fmt.Errorf("invalid regex: %s", err)
	}

	*r = Regexp{*regexpValue}

	return nil
}
