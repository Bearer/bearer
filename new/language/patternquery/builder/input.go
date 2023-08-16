package builder

import (
	"errors"
	"fmt"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/pkg/util/set"
)

func processInput(patternImplementation implementation.Pattern, input string) ([]byte, *InputParams, error) {
	inputWithoutVariables, variables, err := patternImplementation.ExtractVariables(input)
	if err != nil {
		return nil, nil, fmt.Errorf("error processing variables: %s", err)
	}

	inputWithoutVariablesBytes := []byte(inputWithoutVariables)
	matchNodePositions := patternImplementation.FindMatchNode(inputWithoutVariablesBytes)
	inputWithoutMatchNode := stripPositions(inputWithoutVariablesBytes, matchNodePositions)
	matchNodeOffset := 0

	if len(matchNodePositions) > 1 {
		return nil, nil, errors.New("pattern must only contain a single match node")
	}

	if len(matchNodePositions) == 1 {
		matchNodeOffset = matchNodePositions[0][0]
	}

	unanchoredPositions := patternImplementation.FindUnanchoredPoints(inputWithoutMatchNode)
	inputWithoutUnanchored := stripPositions(inputWithoutMatchNode, unanchoredPositions)

	unanchoredOffsets := make([]int, len(unanchoredPositions))
	for i, position := range unanchoredPositions {
		unanchoredOffsets[i] = adjustForPositions(position[0], unanchoredPositions[:i])
	}

	variableNames := set.New[string]()
	for _, variable := range variables {
		variableNames.Add(variable.Name)
	}

	return inputWithoutUnanchored, &InputParams{
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
