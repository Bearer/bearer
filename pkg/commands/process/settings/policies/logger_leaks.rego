package bearer.logger_leaks

import future.keywords

default level := "none"


locations[location] {
    some detector in input.risks
    detector.detector_id == "detect_ruby_logger"
    location = detector.data_types[_].locations[_]
}

level = "warning" if {
    count(locations) > 0
}

