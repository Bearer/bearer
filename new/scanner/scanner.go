package scanner

import (
	"github.com/bearer/curio/new/detector/composition/ruby"
	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/pkg/commands/process/settings"
)

var detectorSet types.DetectorSet

func Setup(config map[string]settings.Rule) (err error) {
	detectorSet, err = ruby.New()
	return err
}

func Detect(file string) (err error) {

}
