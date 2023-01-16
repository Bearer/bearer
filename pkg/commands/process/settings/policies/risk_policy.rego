package bearer.risk_policy

import data.bearer.common

import future.keywords

contains(arr, elem) if {
	# arr # ensure array is defined
	arr[_] = elem
}

policy_failure contains item if {
	input.rule.trigger == "local"
	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id

	data_type = detector.data_types[_]

	not contains(input.rule.exclude_data_types, data_type.name)

	location = data_type.locations[_]
	item := {
		"category_groups": data.bearer.common.groups_for_datatype(data_type),
		"filename": location.filename,
		"line_number": location.line_number,
		"parent_line_number": location.parent.line_number,
		"parent_content": location.parent.content
	}
}

policy_failure contains item if {
	input.rule.trigger == "local"
	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id

	data_type = detector.data_types[_]

	not input.rule.exclude_data_types

	location = data_type.locations[_]
	item := {
		"category_groups": data.bearer.common.groups_for_datatype(data_type),
		"filename": location.filename,
		"line_number": location.line_number,
		"parent_line_number": location.parent.line_number,
		"parent_content": location.parent.content
	}
}

policy_failure contains item if {
	input.rule.trigger == "global"
	not input.rule.omit_parent == false
	not input.rule.omit_parent_content == false

	# FIXME: handle case to use exclude data type with global trigger
	some data_type in input.dataflow.data_types

	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id

	location = detector.locations[_]
	item := {
		"category_groups": data.bearer.common.groups_for_datatypes(input.dataflow.data_types),
		"filename": location.filename,
		"line_number": location.line_number,
		"parent_line_number": location.parent.line_number,
		"parent_content": location.parent.content
	}
}

policy_failure contains item if {
	input.rule.trigger == "global"
	input.rule.omit_parent == true

	# FIXME: handle case to use exclude data type with global trigger
	some data_type in input.dataflow.data_types

	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id

	location = detector.locations[_]
	item := {
		"category_groups": data.bearer.common.groups_for_datatypes(input.dataflow.data_types),
		"filename": location.filename,
		"line_number": location.line_number
	}
}

policy_failure contains item if {
	input.rule.trigger == "global"
	input.rule.omit_parent_content == true

	# FIXME: handle case to use exclude data type with global trigger
	some data_type in input.dataflow.data_types

	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id

	location = detector.locations[_]
	item := {
		"category_groups": data.bearer.common.groups_for_datatypes(input.dataflow.data_types),
		"filename": location.filename,
		"line_number": location.line_number,
		"parent_line_number": location.parent.line_number,
		"parent_content": location.content
	}
}