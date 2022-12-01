package bearer.insecure_smtp

import future.keywords

medium[item] {
    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    location = detector.locations[_]
    item := {
        "category_group": "Insecure communication",
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": location.line_number,
        "parent_content": location.content
    }
}
