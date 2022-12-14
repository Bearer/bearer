package bearer.weak_password_encryption

import data.bearer.common
import data.bearer.encryption_common

import future.keywords

# openssl pkey rsa encryption
policy_failure contains item if {
    some detector in input.dataflow.risks
    detector.detector_id in ["ruby_openssl_pkey_rsa_method_call", "ruby_openssl_pkey_method_call"]

    data_type = detector.data_types[_]
    data_type.uuid == data.bearer.encryption_common.password_uuid

    location = data_type.locations[_]
    location.filename in data.bearer.encryption_common.openssl_pkey_rsa_files

    item := {
        "category_groups": data.bearer.common.groups_for_datatype(data_type),
        "severity": data.bearer.common.severity_of_datatype(data_type),
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": location.parent.line_number,
        "parent_content": location.parent.content
    }
}

# openssl pkey dsa encryption
policy_failure contains item if {
    some detector in input.dataflow.risks
    detector.detector_id in ["ruby_openssl_pkey_dsa_method_call", "ruby_openssl_pkey_method_call"]

    data_type = detector.data_types[_]
    data_type.uuid == data.bearer.encryption_common.password_uuid

    location = data_type.locations[_]
    location.filename in data.bearer.encryption_common.openssl_pkey_dsa_files

    item := {
        "category_groups": data.bearer.common.groups_for_datatype(data_type),
        "severity": data.bearer.common.severity_of_datatype(data_type),
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": location.parent.line_number,
        "parent_content": location.parent.content
    }
}

# blowfish encryption
policy_failure contains item if {
    some detector in input.dataflow.risks
    detector.detector_id == "ruby_blowfish_method_call"

    data_type = detector.data_types[_]
    data_type.uuid == data.bearer.encryption_common.password_uuid

    location = data_type.locations[_]
    location.filename in data.bearer.encryption_common.blowfish_files

    item := {
        "category_groups": data.bearer.common.groups_for_datatype(data_type),
        "severity": data.bearer.common.severity_of_datatype(data_type),
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": location.parent.line_number,
        "parent_content": location.parent.content
    }
}

# rc4 encryption
policy_failure contains item if {
    some detector in input.dataflow.risks
    detector.detector_id == "encrypt_method_call"

    data_type = detector.data_types[_]
    data_type.uuid == data.bearer.encryption_common.password_uuid

    location = data_type.locations[_]
    location.filename in data.bearer.encryption_common.rc4_files

    item := {
        "category_groups": data.bearer.common.groups_for_datatype(data_type),
        "severity": data.bearer.common.severity_of_datatype(data_type),
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": location.parent.line_number,
        "parent_content": location.parent.content
    }
}

policy_failure contains item if {
    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    data_type = detector.data_types[_]
    data_type.uuid == data.bearer.encryption_common.password_uuid

    location = data_type.locations[_]
    item := {
        "category_groups": data.bearer.common.groups_for_datatype(data_type),
        "severity": data.bearer.common.severity_of_datatype(data_type),
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": location.parent.line_number,
        "parent_content": location.parent.content
    }
}
