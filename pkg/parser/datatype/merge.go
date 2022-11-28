package datatype

import (
	"github.com/bearer/curio/pkg/report/schema/datatype"
)

func MergeDatatypesByPropertyNames(source datatype.DataTypable, enrichment datatype.DataTypable) {
	for nodeID, enrichmentDatatype := range enrichment.GetProperties() {
		sourceProperties := source.GetProperties()

		if sourceProperties == nil {
			source.CreateProperties()
			sourceProperties[nodeID] = enrichmentDatatype
			continue
		}

		var sourceChild datatype.DataTypable

		for _, sourceProperty := range sourceProperties {
			if sourceProperty.GetName() == enrichmentDatatype.GetName() {
				sourceChild = sourceProperty
				break
			}
		}

		if sourceChild != nil {
			MergeDatatypesByPropertyNames(sourceChild, enrichmentDatatype)
			continue
		}

		sourceProperties[nodeID] = enrichmentDatatype
	}
}
