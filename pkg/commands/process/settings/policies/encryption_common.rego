package bearer.encryption_common

import future.keywords

password_uuid := "02bb0d3a-2c8c-4842-be1c-c057f0dccd63"

files_with_password_encryption_calls contains item if {
  some detector in input.dataflow.risks
  detector.detector_id in [
    "ruby_openssl_pkey_rsa_method_call",
    "ruby_openssl_pkey_method_call",
    "ruby_blowfish_method_call",
    "encrypt_method_call"
  ]

  data_type = detector.data_types[_]
  data_type.uuid == password_uuid

  location = detector.locations[_]

  item := location.filename
}

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