package entropy_test

import (
	"fmt"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/pkg/util/entropy"
)

func TestShannon(t *testing.T) {
	examples := []string{
		"one",
		"password",
		"secret_key",
		"Au+u1hvsvJeEXxky",
	}

	results := make([]string, len(examples))
	for i, example := range examples {
		results[i] = fmt.Sprintf("%.2f", entropy.Shannon(example))
	}

	cupaloy.SnapshotT(t, results)
}
