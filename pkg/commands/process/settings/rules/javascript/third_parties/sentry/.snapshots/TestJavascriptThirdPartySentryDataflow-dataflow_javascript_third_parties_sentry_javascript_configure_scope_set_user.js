data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/sentry/testdata/javascript_configure_scope_set_user.js
              line_number: 2
              field_name: email
              object_name: setUser
              subject_name: User
risks:
    - detector_id: javascript_third_parties_sentry_scope
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/sentry/testdata/javascript_configure_scope_set_user.js
          line_number: 1
          parent:
            line_number: 1
            content: scope
          content: |
            Sentry.configureScope(($<!>$<_:identifier>) => {})
components: []


--

