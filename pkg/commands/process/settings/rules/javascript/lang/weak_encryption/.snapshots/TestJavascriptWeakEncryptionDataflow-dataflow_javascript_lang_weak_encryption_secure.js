data_types:
    - name: Passwords
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/secure.js
              line_number: 4
              field_name: password
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/secure.js
              line_number: 5
              field_name: password
              object_name: user
              subject_name: User
components: []


--

