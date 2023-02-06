package bearer.privacy_report

import data.bearer.common

import future.keywords

items contains item if {
  some data_type in input.dataflow.data_types
  some detector in data_type.detectors
  some location in detector.locations

  item := {
    "name": data_type.name,
    "subject_name": location.subject_name,
    "line_number": location.line_number
  }
}
