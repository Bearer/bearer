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