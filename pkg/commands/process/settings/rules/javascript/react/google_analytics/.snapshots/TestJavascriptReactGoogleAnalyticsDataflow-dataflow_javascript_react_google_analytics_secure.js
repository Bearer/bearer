data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/react/google_analytics/testdata/secure.js
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
components: []


--

