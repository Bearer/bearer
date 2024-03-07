package bearer.common

import rego.v1

build_item(location) := {
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
		"end": location.end_line_number,
		"column": {
			"start": location.start_column_number,
			"end": location.end_column_number,
		},
	},
	"line_number": location.start_line_number,
	"detailed_context": location.presence_matches[0].name,
} if {
	# FIXME: This is only for secret detections. Should be more explicit
	input.rule.has_detailed_context == true
}

cat_groups := data.bearer.common.groups_for_datatypes(input.dataflow.data_types) if {
	some data_type in input.dataflow.data_types
}

cat_groups := set() if {
	not input.dataflow.data_types
}

build_local_item(location, data_type) := {
	"is_local": true,
	"category_groups": groups_for_datatype(data_type),
	"data_type": {
		"category_uuid": data_type.category_uuid,
		"name": data_type.name,
	},
	"full_filename": location.full_filename,
	"filename": location.filename,
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
		"end": location.end_line_number,
		"column": {
			"start": location.start_column_number,
			"end": location.end_column_number,
		},
	},
	"line_number": location.source.start_line_number,
} if {
	not input.rule.has_detailed_context == true
}

build_item(location) := {
	"category_groups": cat_groups,
	"full_filename": location.full_filename,
	"filename": location.filename,
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
		"end": location.end_line_number,
		"column": {
			"start": location.start_column_number,
			"end": location.end_column_number,
		},
	},
	"line_number": location.start_line_number,
} if {
	not input.rule.has_detailed_context == true
}

global_data_types contains data_type if {
	not input.rule.only_data_types
	not input.rule.skip_data_types

	some data_type in input.dataflow.data_types
}

global_data_types contains data_type if {
	not input.rule.only_data_types

	some data_type in input.dataflow.data_types
	not contains(input.rule.skip_data_types, data_type.name)
}

global_data_types contains data_type if {
	not input.rule.skip_data_types

	some data_type in input.dataflow.data_types
	contains(input.rule.only_data_types, data_type.name)
}

groups_for_datatype(data_type) := x if {
	some category in input.data_categories
	category.uuid == data_type.category_uuid

	x := {name | name := category.groups[_].name}
}

groups_for_datatypes(data_types) := groups if {
	groups := {name | name := groups_for_datatype(data_types[_])[_]}
}
