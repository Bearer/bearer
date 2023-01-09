package bearer.secret_leak

import data.bearer.common

import future.keywords

policy_failure contains item if {
	some detector in input.dataflow.risks
	detector.detector_id == "gitleaks"

	location = detector.locations[_]
	item := {
		"severity": "critical",
		"filename": location.filename,
		"parent_line_number": location.line_number,
		"line_number": location.line_number,
		"detailed_context": location.content,
	}
}
