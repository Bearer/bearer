package bearer.ssl_certificate_verification_disabled

import data.bearer.common

import future.keywords

policy_failure contains item if {
  some data_type in input.dataflow.data_types

  some detector in input.dataflow.risks
  detector.detector_id == input.policy_id
  location = detector.locations[_]

  item = {
    "category_groups": data.bearer.common.groups_for_datatype(data_type),
    "severity": "medium",
    "filename": location.filename,
    "line_number": location.line_number,
    "parent_line_number": location.parent.line_number,
    "parent_content": location.parent.content
  }
}
