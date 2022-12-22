package curio.detectors.rego_properties

import future.keywords.contains
import future.keywords.if

pair_query := curio.language.compile_sitter_query(`(pair key: (hash_key_symbol) @key value: (_) @value) @root`)

detections_at contains detections if {
	match := curio.query.match_once_at(pair_query, input)
	match != null

	name := curio.node.content(match.key)

	detections := { "match_node": input, "data": { "name": name }}
}
