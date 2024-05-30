package datatype

import (
	"github.com/bearer/bearer/pkg/report/schema/datatype"
)

// MergeDatatypesByPropertyNames iterates trough source properties and adds missing datatypes from enrichment
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

// CloneDeepestProperties deep clones source and iterates trough various levels
// of same name source and enrichement properties adding enrichment properties
// to the edge of source only when source edge doesn't have any properties
func CloneDeepestProperties(source datatype.DataTypable, enrichment datatype.DataTypable) datatype.DataTypable {
	clonedSource := source.Clone()

	sourceProperties := clonedSource.GetProperties()
	if len(sourceProperties) == 0 {
		clonedSource.CreateProperties()
		for nodeID, enenrichmentChild := range enrichment.GetProperties() {
			clonedSource.SetProperty(nodeID, enenrichmentChild)
		}
		return clonedSource
	}

	for enrichementNodeID, enrichmentDatatype := range enrichment.GetProperties() {
		var sourceChild datatype.DataTypable
		for _, sourceProperty := range sourceProperties {
			if sourceProperty.GetName() == enrichmentDatatype.GetName() {
				sourceChild = sourceProperty
				break
			}
		}

		if sourceChild != nil {
			clonedSource.SetProperty(enrichementNodeID, CloneDeepestProperties(sourceChild, enrichmentDatatype))
			continue
		}
	}

	return clonedSource
}
