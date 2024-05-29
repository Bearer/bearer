package entropy

import "math"

func Shannon(value string) float64 {
	if value == "" {
		return 0
	}

	counts := make(map[rune]int)
	for _, character := range value {
		counts[character]++
	}

	length := float64(len(value))
	var negativeEntropy float64

	for _, count := range counts {
		characterProbability := float64(count) / length
		negativeEntropy += characterProbability * math.Log2(characterProbability)
	}

	return -negativeEntropy
}
