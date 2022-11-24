package bearer.encrypted_tanker

import future.keywords


default encrypted := false

encrypted = true {
    startswith(lower(input.target.value.field_name),  "tanker_encrypted")
}

verified_by := [
    {
        "detector": "tanker_encrypted",
    }
]