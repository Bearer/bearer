package bearer.common

import future.keywords

build_item(location) := {
	"filename": location.filename,
	"parent_line_number": location.parent.line_number,
	"parent_content": location.parent.content,
	"line_number": location.line_number,
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
	"filename": location.filename,
	"line_number": location.line_number,
	"parent_line_number": location.parent.line_number,
	"parent_content": location.parent.content,
	"datatype_name": data_type.name
} if {
	not input.rule.has_detailed_context == true
}

build_item(location) := {
	"category_groups": cat_groups,
	"filename": location.filename,
	"line_number": location.line_number,
	"parent_line_number": location.parent.line_number,
	"parent_content": location.parent.content,
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
