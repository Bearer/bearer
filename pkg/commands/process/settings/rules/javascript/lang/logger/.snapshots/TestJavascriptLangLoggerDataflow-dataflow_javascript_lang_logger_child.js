risks:
    - detector_id: javascript_lang_logger
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/logger/testdata/child.js
              line_number: 3
              parent:
                line_number: 7
                content: logger.child(ctx)
              field_name: email
              object_name: user
              subject_name: User
    - detector_id: child_logger
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/lang/logger/testdata/child.js
          line_number: 7
          parent:
            line_number: 7
            content: logger.child(ctx)
          content: |
            $<LOG>.child()
components: []


--

