data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/sentry/testdata/javascript_capture_exception.js
              line_number: 2
              field_name: email
              object_name: current_user
              subject_name: User
risks:
    - detector_id: javascript_third_parties_sentry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/sentry/testdata/javascript_capture_exception.js
              line_number: 2
              parent:
                line_number: 1
                content: |-
                    Sentry.captureException(
                      new Error(`user ${current_user.email} couldn't log in!`)
                    )
              field_name: email
              object_name: current_user
              subject_name: User
components: []


--

