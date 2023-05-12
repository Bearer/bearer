package bearer.rules

import future.keywords.contains
import future.keywords.if
import future.keywords.in

pattern_match contains [node, pattern_id, variables] if {
  some match in input.pattern_matches

  node := match.node
  pattern_id := match.pattern_id
  variables := match.variables
}

