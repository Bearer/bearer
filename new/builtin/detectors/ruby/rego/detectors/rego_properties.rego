package curio.detectors.rego_properties

import future.keywords.contains
import future.keywords.if

pair_query := curio.language.compile_sitter_query(`(pair key: (hash_key_symbol) @key value: (_) @value) @root`)

detections_at(node) := detections if {
	match := curio.query.match_once_at(pair_query, node)
	match != null

	name := curio.node.content(match.key)

	detections := { "match_node": node, "data": { "name": name }}
}

detections_at_input contains detections if {
	detections := detections_at(input)
}
