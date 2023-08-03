package builder

import (
	"errors"
	"fmt"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/pkg/util/set"
)

func processInput(langImplementation implementation.Implementation, input string) (string, *InputParams, error) {
	inputWithoutVariables, variables, err := langImplementation.ExtractPatternVariables(input)
	if err != nil {
		return "", nil, fmt.Errorf("error processing variables: %s", err)
	}

	inputWithoutVariablesBytes := []byte(inputWithoutVariables)
	matchNodePositions := langImplementation.FindPatternMatchNode(inputWithoutVariablesBytes)
	inputWithoutMatchNode := stripPositions(inputWithoutVariablesBytes, matchNodePositions)
	matchNodeOffset := 0

	if len(matchNodePositions) > 1 {
		return "", nil, errors.New("pattern must only contain a single match node")
	}

	if len(matchNodePositions) == 1 {
		matchNodeOffset = matchNodePositions[0][0]
	}

	unanchoredPositions := langImplementation.FindPatternUnanchoredPoints(inputWithoutMatchNode)
	inputWithoutUnanchored := stripPositions(inputWithoutMatchNode, unanchoredPositions)

	unanchoredOffsets := make([]int, len(unanchoredPositions))
	for i, position := range unanchoredPositions {
		unanchoredOffsets[i] = adjustForPositions(position[0], unanchoredPositions[:i])
	}

	variableNames := set.New[string]()
	for _, variable := range variables {
		variableNames.Add(variable.Name)
	}

	return string(inputWithoutUnanchored), &InputParams{
		Variables:         variables,
		VariableNames:     variableNames.Items(),
		MatchNodeOffset:   adjustForPositions(matchNodeOffset, unanchoredPositions),
		UnanchoredOffsets: unanchoredOffsets,
	}, nil
}

func stripPositions(input []byte, positions [][]int) []byte {
	offset := 0
	var result []byte

	for _, position := range positions {
		result = append(result, input[offset:position[0]]...)
		offset = position[1]
	}

	return append(result, input[offset:]...)
}

func adjustForPositions(offset int, positions [][]int) int {
	result := offset

	for _, position := range positions {
		if position[0] >= offset {
			break
		}

		result -= (position[1] - position[0])
	}

	return result
}
