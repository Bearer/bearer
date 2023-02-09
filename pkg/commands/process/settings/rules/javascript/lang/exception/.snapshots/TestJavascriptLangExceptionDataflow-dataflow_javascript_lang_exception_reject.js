data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/exception/testdata/reject.js
              line_number: 5
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/lang/exception/testdata/reject.js
              line_number: 14
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: javascript_lang_exception
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/exception/testdata/reject.js
              line_number: 5
              parent:
                line_number: 7
                content: reject("Error with user " + user)
              field_name: email
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/lang/exception/testdata/reject.js
              line_number: 14
              parent:
                line_number: 16
                content: reject("Error with user " + user)
              field_name: email
              object_name: current_user
              subject_name: User
    - detector_id: javascript_lang_new_promise_init
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/lang/exception/testdata/reject.js
          line_number: 2
          parent:
            line_number: 2
            content: reject
          content: |
            new Promise(function ($<_>, $<!>$<_>) {})
        - filename: pkg/commands/process/settings/rules/javascript/lang/exception/testdata/reject.js
          line_number: 11
          parent:
            line_number: 11
            content: reject
          content: |
            new Promise(($<_>, $<!>$<_>) => {})
components: []


--

