package bearer.risk_policy

import rego.v1

import data.bearer.common

array_contains(arr, elem) if {
	arr[_] = elem
}

has_key(x, k) if _ = x[k]

# - presence of pattern & data types required
global_failures contains detector if {
	input.rule.trigger.match_on == "presence"
	input.rule.trigger.data_types_required

	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id

	some data_type in data.bearer.common.global_data_types
}

presence_failures contains detector if {
	input.rule.trigger.match_on == "presence"
	not input.rule.trigger.data_types_required

	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id

	not input.rule.dependency_check
}

# - presence of pattern & data types not required
presence_failures contains detector if {
	input.rule.trigger.match_on == "presence"
	not input.rule.trigger.data_types_required

	some detector in input.dataflow.risks
	detector.detector_id == input.rule.id

	input.rule.dependency_check
	dependency := input.dataflow.dependencies[_]
	dependency.filename == input.rule.dependency.filename
	dependency.name == input.rule.dependency.name

	version_check := semver.compare(dependency.version, input.rule.dependency.min_version)
	version_check <= 0
}

# # Build policy failures
policy_failure contains item if {
	input.rule.trigger.match_on == "absence"

	count(input.rule.trigger.required_detections) == count({ required_detection |
		required_detection := input.rule.trigger.required_detections[_]
		some y in input.dataflow.risks
		y.detector_id == required_detection
	})

	some detector in input.dataflow.risks
	detector.detector_id == input.rule.trigger.required_detections[0]

	some init_location in detector.locations

	x := {other | other := input.dataflow.risks[_]; other.detector_id == input.rule.id}
	count(x) == 0
	item := data.bearer.common.build_item(init_location)
}

policy_failure contains item if {
	input.rule.trigger.match_on == "absence"
	count(input.rule.trigger.required_detections) == count({ required_detection |
		required_detection := input.rule.trigger.required_detections[_]
		some x in input.dataflow.risks
		x.detector_id == required_detection
	})

	some detector in input.dataflow.risks
	detector.detector_id == input.rule.trigger.required_detections[0]

	some init_location in detector.locations
	some other_detector in input.dataflow.risks

	other_detector.detector_id == input.rule.id

	x := {other_location | other_location := other_detector.locations[_]; init_location.filename == other_location.filename}
	count(x) == 0

	item := data.bearer.common.build_item(init_location)
}

# - data types detected within pattern ($<DATA_TYPE>)
policy_failure contains item if {
	not input.rule.skip_data_types
	not input.rule.only_data_types

	some detector in presence_failures
	some data_type_location in detector.locations
	some data_type in data_type_location.data_types

	item := data.bearer.common.build_local_item(data_type_location, data_type)
}

policy_failure contains item if {
	not input.rule.skip_data_types

	some detector in presence_failures
	some data_type_location in detector.locations
	some data_type in data_type_location.data_types
	array_contains(input.rule.only_data_types, data_type.name)

	item := data.bearer.common.build_local_item(data_type_location, data_type)
}

policy_failure contains item if {
	not input.rule.only_data_types

	some detector in presence_failures
	some data_type_location in detector.locations
	some data_type in data_type_location.data_types
	not array_contains(input.rule.skip_data_types, data_type.name)

	item := data.bearer.common.build_local_item(data_type_location, data_type)
}

# end datatyped detection

# policies for global failures
policy_failure contains item if {
	some detector in global_failures

	location = detector.locations[_]
	item := data.bearer.common.build_item(location)
}

# policies for presence failures
policy_failure contains item if {
	some detector in presence_failures
	some location in detector.locations
	not has_key(location, "data_types")

	item := data.bearer.common.build_item(location)
}

policy_failure contains item if {
	input.rule.trigger.match_on == "stored_data_types"

	array_contains(input.rule.languages, input.dataflow.data_types[_].detectors[_].name)
	data_type = input.dataflow.data_types[_]
	not array_contains(input.rule.skip_data_types, data_type.name)

	some detector in data_type.detectors

	array_contains(input.rule.detectors, detector.name)

	location = detector.locations[_]
	count(input.rule.auto_encrypt_prefix) != 0

	not location.encrypted == true

	item := {
		"category_groups": data.bearer.common.groups_for_datatype(data_type),
		"filename": location.filename,
		"full_filename": location.full_filename,
		"sink": {
			"start": location.source.start_line_number,
			"end": location.source.end_line_number,
			"column": {
				"start": location.source.start_column_number,
				"end": location.source.end_column_number,
			},
		},
		"source": {
			"start": location.start_line_number,
			"end": location.start_line_number,
			"column": {
				"start": location.start_column_number,
				"end": location.end_column_number,
			},
		},
		"line_number": location.start_line_number,
	}
}

# used by inventory report
local_rule_failure contains item if {
	some detector in presence_failures
	some location in detector.locations
	some data_type in location.data_types
	some schema in data_type.schemas

	item := {
		"name": data_type.name,
		"category_groups": data.bearer.common.groups_for_datatype(data_type),
		"subject_name": schema.subject_name,
		"line_number": location.start_line_number,
		"rule_id": input.rule.id,
		"third_party": input.rule.associated_recipe,
	}
}
