package bearer.inventory_report

import data.bearer.common

import future.keywords

report_items contains item if {
  some data_type in input.dataflow.data_types
  some detector in data_type.detectors
  some location in detector.locations

  item := {
    "name": data_type.name,
    "object_name": location.object_name,
    "line_number": location.line_number
  }
}
