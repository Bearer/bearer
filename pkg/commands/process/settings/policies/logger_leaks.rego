package bearer.logger_leaks

import future.keywords

default warning := false

warning if {
some detector in input.risks
detector.detector_id == "detect_ruby_logger"
}