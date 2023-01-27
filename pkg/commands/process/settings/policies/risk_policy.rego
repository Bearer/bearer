package bearer.risk_policy

import data.bearer.common

import future.keywords

contains(arr, elem) if {
	arr[_] = elem
}

local_failures contains detector if {
	input.rule.trigger == "local"
	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id
}

global_failures contains detector if {
	input.rule.trigger == "global"
	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id

	some data_type in data.bearer.common.global_data_types
}

presence_failures contains detector if {
	input.rule.trigger == "presence"
	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id
}

local_data_types contains data_type if {
	not input.rule.skip_data_types
	not input.rule.only_data_types

	some detector in local_failures
	data_type = detector.data_types[_]
}

local_data_types contains data_type if {
	not input.rule.only_data_types

	some detector in local_failures
	data_type = detector.data_types[_]
	not contains(input.rule.skip_data_types, data_type.name)
}

local_data_types contains data_type if {
	not input.rule.skip_data_types

	some detector in local_failures
	data_type = detector.data_types[_]
	contains(input.rule.only_data_types, data_type.name)
}

# Build policy failures

policy_failure contains item if {
	some data_type in local_data_types

	location = data_type.locations[_]
	item := data.bearer.common.build_item(location)
}

policy_failure contains item if {
	some detector in global_failures

	location = detector.locations[_]
	item := data.bearer.common.build_item(location)
}

policy_failure contains item if {
	some detector in presence_failures

	location = detector.locations[_]
	item := data.bearer.common.build_item(location)
}

policy_failure contains item if {
	input.rule.trigger == "stored_data_types"

	data_type = input.dataflow.data_types[_]
	not contains(input.rule.skip_data_types, data_type.name)

	some detector in data_type.detectors

	contains(input.rule.detectors, detector.name)

	location = detector.locations[_]
	count(input.rule.auto_encrypt_prefix) != 0

	not location.encrypted == true

	item := {
		"category_groups": data.bearer.common.groups_for_datatype(data_type),
		"filename": location.filename,
		"line_number": location.line_number,
		"parent_line_number": location.parent.line_number,
		"parent_content": location.parent.content,
	}
}

# used by inventory report
local_rule_failure contains item if {
  some detector in local_failures
	data_type = detector.data_types[_]

  location = data_type.locations[_]
	item := {
    "name": data_type.name,
		"category_groups": data.bearer.common.groups_for_datatype(data_type),
    "object_name": location.object_name,
    "line_number": location.line_number,
    "rule_id": input.rule.id
  }
}