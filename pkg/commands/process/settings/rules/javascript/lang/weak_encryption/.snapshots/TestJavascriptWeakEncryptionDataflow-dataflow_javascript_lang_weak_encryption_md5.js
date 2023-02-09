data_types:
    - name: Passwords
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/md5.js
              line_number: 4
              field_name: password
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/md5.js
              line_number: 5
              field_name: password
              object_name: user
              subject_name: User
risks:
    - detector_id: javascript_weak_encryption
      data_types:
        - name: Passwords
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/md5.js
              line_number: 4
              parent:
                line_number: 4
                content: crypto.createHmac("md5", key).update(user.password)
              field_name: password
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/md5.js
              line_number: 5
              parent:
                line_number: 5
                content: crypto.createHash("md5").update(user.password)
              field_name: password
              object_name: user
              subject_name: User
    - detector_id: create_hash
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/md5.js
          line_number: 4
          parent:
            line_number: 4
            content: crypto.createHmac("md5", key)
          content: |
            crypto.$<METHOD>($<ALGORITHM>$<...>)
        - filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/md5.js
          line_number: 5
          parent:
            line_number: 5
            content: crypto.createHash("md5")
          content: |
            crypto.$<METHOD>($<ALGORITHM>$<...>)
components: []


--

