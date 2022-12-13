package bearer.encryption_common

import future.keywords

rc4_files contains item if {
  some detector in input.dataflow.risks
  detector.detector_id == "initialize_ruby_rc4_encryption"

  location = detector.locations[_]

  item := location.filename
}

blowfish_files contains item if {
  some detector in input.dataflow.risks
  detector.detector_id == "initialize_ruby_blowfish_encryption"

  location = detector.locations[_]

  item := location.filename
}

openssl_pkey_rsa_files contains item if {
  some detector in input.dataflow.risks
  detector.detector_id == "initialize_ruby_openssl_pkey_rsa_encryption"

  location = detector.locations[_]

  item := location.filename
}

openssl_pkey_dsa_files contains item if {
  some detector in input.dataflow.risks
  detector.detector_id == "initialize_ruby_openssl_pkey_dsa_encryption"

  location = detector.locations[_]

  item := location.filename
}