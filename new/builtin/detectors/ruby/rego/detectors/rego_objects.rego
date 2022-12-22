package curio.detectors.rego_objects

import future.keywords.contains
import future.keywords.if


# { first_name: ..., ... }
hash_pair_query := curio.language.compile_sitter_query(`(hash (pair) @pair) @root`)

# gather properties ( { first_name: "", last_name: "" } )
detections_at contains detections if {
	results := curio.query.match_at(hash_pair_query, input)
	count(results) != 0

  properties := [property |
    property := curio.evaluator.detections_at(results[_].pair, "rego_properties")[_]
  ]

  detections := [{ "match_node": input, "data": { "properties": properties }}]
}


# user = <object>
assignment_query := curio.language.compile_sitter_query(`(assignment left: (identifier) @left right: (_) @right) @root`)

# name assigned objects ( user = ... )
detections_at contains detections if {
	result := curio.query.match_once_at(assignment_query, input)
  result

	objects := curio.evaluator.detections_at(result.right, "rego_objects")
  count(objects) != 0

  detections := [detection |
    object := objects[_]
    detection := wrap_assigned_object(input, result, object)
  ]
}

wrap_assigned_object(node, result, object) := detection {
  data := object.data
  data.name == ""

  detection := {
    "match_node": node,
    "data": {
      "name": curio.node.content(result.left),
      "properties": data.properties
    }
  }
}

wrap_assigned_object(node, result, object) := detection {
  data := object.data
  data.name != ""

  detection := {
    "match_node": node,
    "data": {
      "name": curio.node.content(result.left),
      "properties": [
        {
          "match_node": object.match_node,
          "data": { "name": data.name }
        }
      ]
    }
  }
}


# { user: <object> }
parent_pair_query := curio.language.compile_sitter_query(`(pair key: (hash_key_symbol) @key value: (_) @value) @root`)

# name parent-pair objects ( { user: ... } )
detections_at contains detections if {
	result := curio.query.match_once_at(parent_pair_query, input)
  result

	objects := curio.evaluator.detections_at(result.value, "rego_objects")

  detections := [detection |
    object := objects[_]

    detection := {
      "match_node": input,
      "data": {
        "name": curio.node.content(result.key),
        "properties": object.data.properties
      }
    }
  ]
}
