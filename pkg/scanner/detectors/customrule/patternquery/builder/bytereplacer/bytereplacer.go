package bytereplacer

import (
	"bytes"
	"fmt"
)

type Replacer struct {
	offsetDelta
	result *Result
}

type offsetDelta struct {
	originalOffset,
	delta int
}

type Result struct {
	value        []byte
	offsetDeltas []offsetDelta
	changed      bool
}

func New(original []byte) *Replacer {
	value := make([]byte, len(original))
	copy(value, original)

	return &Replacer{result: &Result{value: value}}
}

func (replacer *Replacer) Replace(originalStart, originalEnd int, newValue []byte) error {
	if originalStart < replacer.originalOffset {
		return fmt.Errorf(
			"replacements must be made in sequential order. last offset %d, replacement start %d",
			replacer.originalOffset,
			originalStart,
		)
	}

	if bytes.Equal(replacer.result.value[originalStart:originalEnd], newValue) {
		return nil
	}

	replacer.result.changed = true

	currentLength := len(replacer.result.value)
	suffix := replacer.result.value[replacer.delta+originalEnd : currentLength : currentLength]

	replacer.result.value = append(
		replacer.result.value[:replacer.delta+originalStart],
		append(newValue, suffix...)...,
	)

	delta := len(newValue) - originalEnd + originalStart

	replacer.originalOffset = originalEnd
	replacer.delta += delta

	replacer.result.offsetDeltas = append(replacer.result.offsetDeltas, offsetDelta{
		originalOffset: originalEnd,
		delta:          delta,
	})

	return nil
}

func (replacer *Replacer) Done() *Result {
	return replacer.result
}

func (result *Result) Changed() bool {
	return result.changed
}

func (result *Result) Value() []byte {
	return result.value
}

func (result *Result) Translate(originalOffset int) int {
	delta := 0

	for _, offsetDelta := range result.offsetDeltas {
		if offsetDelta.originalOffset > originalOffset {
			break
		}

		delta += offsetDelta.delta
	}

	return delta + originalOffset
}
