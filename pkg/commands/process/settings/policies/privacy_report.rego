package bearer.privacy_report

import rego.v1

import data.bearer.common

array_contains(arr, elem) if {
	arr[_] = elem
}

data_types_with_subject contains item if {
	some data_type in input.dataflow.data_types
	some detector in data_type.detectors
	some location in detector.locations

	location.subject_name

	item = data_type.name
}

items contains item if {
	some data_type in input.dataflow.data_types
	some detector in data_type.detectors
	some location in detector.locations

	location.subject_name

	item := {
		"name": data_type.name,
		"subject_name": location.subject_name,
		"line_number": location.start_line_number,
	}
}

items contains item if {
	some data_type in input.dataflow.data_types
	some detector in data_type.detectors
	some location in detector.locations

	not location.subject_name
	not array_contains(data_types_with_subject, data_type.name)

	item := {
		"name": data_type.name,
		"subject_name": "",
		"line_number": location.start_line_number,
	}
}
