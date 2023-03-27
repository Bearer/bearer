package bearer.risk_policy

import data.bearer.common

import future.keywords

contains(arr, elem) if {
	arr[_] = elem
}

# - presence of pattern & data types required
global_failures contains detector if {
	input.rule.trigger.match_on == "presence"
	input.rule.trigger.data_types_required

	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id

	some data_type in data.bearer.common.global_data_types
}

# - presence of pattern & data types not required
presence_failures contains detector if {
	input.rule.trigger.match_on == "presence"
	not input.rule.trigger.data_types_required

	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id
}

# - data types detected within pattern ($<DATA_TYPE>)
local_data_types contains data_type if {
	not input.rule.skip_data_types
	not input.rule.only_data_types

	some detector in presence_failures
	data_type = detector.locations[_].data_types[_]
}

local_data_types contains data_type if {
	not input.rule.only_data_types

	some detector in presence_failures
	data_type = detector.locations[_].data_types[_]
	not contains(input.rule.skip_data_types, data_type.name)
}

local_data_types contains data_type if {
	not input.rule.skip_data_types

	some detector in presence_failures
	data_type = detector.locations[_].data_types[_]
	contains(input.rule.only_data_types, data_type.name)
}

# Build policy failures
policy_failure contains item if {
	input.rule.trigger.match_on == "absence"
	some detector in input.dataflow.risks

	detector.detector_id == input.rule.trigger.required_detection
	some init_location in detector.locations

	x := {other | other := input.dataflow.risks[_]; other.detector_id == input.rule.id}
	count(x) == 0

	item := data.bearer.common.build_item(init_location)
}

policy_failure contains item if {
	input.rule.trigger.match_on == "absence"
	some detector in input.dataflow.risks

	detector.detector_id == input.rule.trigger.required_detection

	some init_location in detector.locations
	some other_detector in input.dataflow.risks

	other_detector.detector_id == input.rule.id

	x := {other_location | other_location := other_detector.locations[_]; init_location.filename == other_location.filename}
	count(x) == 0

	item := data.bearer.common.build_item(init_location)
}

policy_failure contains item if {
	some data_type in local_data_types

	location = data_type.locations[_]
	item := data.bearer.common.build_local_item(location, data_type)
}

policy_failure contains item if {
	some detector in global_failures

	location = detector.locations[_]
	item := data.bearer.common.build_item(location)
}

policy_failure contains item if {
	some detector in presence_failures
	count(local_data_types) == 0 # detector item already included (through local_data_types)

	location = detector.locations[_]
	item := data.bearer.common.build_item(location)
}

policy_failure contains item if {
	input.rule.trigger.match_on == "stored_data_types"

	contains(input.rule.languages, input.dataflow.data_types[_].detectors[_].name)
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
	some detector in presence_failures
	some location in detector.locations
	some data_type in location.data_types

	# location = data_type.locations[_]
	item := {
		"name": data_type.name,
		"category_groups": data.bearer.common.groups_for_datatype(data_type),
		"subject_name": location.subject_name,
		"line_number": location.line_number,
		"rule_id": input.rule.id,
		"third_party": input.rule.associated_recipe,
	}
}