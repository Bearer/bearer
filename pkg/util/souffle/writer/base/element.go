package base

import (
	"fmt"
	"strconv"
	"strings"
)

type Symbol string
type Unsigned uint32
type Record []Element

type Element interface {
	String() string
	sealedElement()
}

func (Symbol) sealedElement()   {}
func (Unsigned) sealedElement() {}
func (Record) sealedElement()   {}

func (unsigned Unsigned) String() string {
	return fmt.Sprintf("%d", unsigned)
}

func (symbol Symbol) String() string {
	return strconv.Quote(string(symbol))
}

func (record Record) String() string {
	elementsStrings := make([]string, len(record))
	for i, recordElement := range record {
		elementsStrings[i] = ElementString(recordElement)
	}

	return fmt.Sprintf("[%s]", strings.Join(elementsStrings, ", "))
}

func ElementString(element interface{ String() string }) string {
	if element == nil {
		return "nil"
	}

	return element.String()
}
