package schema

import (
	"github.com/bearer/curio/pkg/parser/datatype"
)

type ClassifiedDatatype struct {
	*datatype.DataType
	Properties     map[string]ClassifiedDatatype
	Classification Classification `json:"classification"`
}

type Classification struct {
	Name string
}

type Classifier struct {
	config Config
}

type Config struct {
}

func New(config Config) *Classifier {
	return &Classifier{config: config}
}

func (classifier *Classifier) Classify(data datatype.DataType) (ClassifiedDatatype, error) {
	// todo: implement interface classification (bigbear etc...)
	return ClassifiedDatatype{
		DataType: &datatype.DataType{
			UUID: "1",
		},
		Classification: Classification{
			Name: "personal data",
		},
		Properties: map[string]ClassifiedDatatype{
			"address": {
				Classification: Classification{
					Name: "personal data",
				},
				DataType: &datatype.DataType{
					UUID: "2",
				},
			},
			"age": {
				Classification: Classification{},
			},
		},
	}, nil
}
