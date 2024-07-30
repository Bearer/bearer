package bearer.verifier_policy

import rego.v1

import data.bearer.common

array_contains(arr, elem) if {
	# arr # ensure array is defined
	arr[_] = elem
}

policy_failure contains item if {
	input.rule.trigger == "stored_data_types"

	data_type = input.dataflow.data_types[_]
	not array_contains(input.rule.skip_data_types, data_type.name)

	some detector in data_type.detectors

	array_contains(input.rule.detectors, detector.name)

	location = detector.locations[_]
	count(input.rule.auto_encrypt_prefix) != 0

	not location.encrypted == true
	not startswith(location.field_name, input.rule.auto_encrypt_prefix)

	item := {
		"category_groups": data.bearer.common.groups_for_datatype(data_type),
		"relative_path": location.relative_path,
		"filename": location.filename,
		"source": {
			"start": location.source.start_line_number,
			"end": location.source.end_line_number,
			"content": location.source.content,
			"column": {
				"start": location.source.start_column_number,
				"end": location.source.end_column_number,
			},
		},
		"sink": {
			"start": location.start_line_number,
			"end": location.end_line_number,
			"column": {
				"start": location.start_column_number,
				"end": location.end_column_number,
			},
		},
		"line_number": location.start_line_number,
	}
}

policy_failure contains item if {
	input.rule.trigger == "stored_data_types"

	data_type = input.dataflow.data_types[_]

	not array_contains(input.rule.skip_data_types, data_type.name)
	some detector in data_type.detectors

	array_contains(input.rule.detectors, detector.name)

	location = detector.locations[_]

	location.stored == true
	not location.encrypted == true

	count(input.rule.auto_encrypt_prefix) == 0
	not location.encrypted == true

	item := {
		"category_groups": data.bearer.common.groups_for_datatype(data_type),
		"full_filename": location.full_filename,
		"filename": location.filename,
		"sink": {
			"start": location.source.start_line_number,
			"end": location.source.end_line_number,
			"content": location.source.content,
			"column": {
				"start": location.source.start_column_number,
				"end": location.source.end_column_number,
			},
		},
		"source": {
			"start": location.start_line_number,
			"end": location.end_line_number,
			"column": {
				"start": location.start_column_number,
				"end": location.end_column_number,
			},
		},
		"line_number": location.start_line_number,
	}
}
