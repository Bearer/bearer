package bearer.data_type_policy

import data.bearer.common

import future.keywords

contains(arr, elem) if {
	# arr # ensure array is defined
	arr[_] = elem
}

policy_failure contains item if {
	input.rule.stored == true

	data_type = input.dataflow.data_types[_]

	some detector in data_type.detectors

	contains(array.concat(input.rule.linked_detectors, [input.rule.id]), detector.name)

	location = detector.locations[_]

	count(input.rule.auto_encrypt_prefix) != 0

	# location.encrypted != false

	not startswith(location.field_name, input.rule.auto_encrypt_prefix)

	item := {
		"category_groups": data.bearer.common.groups_for_datatype(data_type),
		"filename": location.filename,
		"line_number": location.line_number,
		"parent_line_number": location.parent.line_number,
		"parent_content": location.parent.content,
	}
}
