data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/datadog/testdata/unsecure.js
              line_number: 3
              field_name: email
              object_name: user
              subject_name: User
    - name: Fullname
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/datadog/testdata/unsecure.js
              line_number: 3
              field_name: name
              object_name: user
              subject_name: User
risks:
    - detector_id: javascript_third_parties_datadog
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/datadog/testdata/unsecure.js
              line_number: 3
              parent:
                line_number: 11
                content: client.event("user", "logged_in", {}, user)
              field_name: email
              object_name: user
              subject_name: User
        - name: Fullname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/datadog/testdata/unsecure.js
              line_number: 3
              parent:
                line_number: 11
                content: client.event("user", "logged_in", {}, user)
              field_name: name
              object_name: user
              subject_name: User
    - detector_id: javascript_third_parties_hotshot_statsd
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/datadog/testdata/unsecure.js
          line_number: 5
          parent:
            line_number: 5
            content: |-
                new StatsD({
                	port: 8020,
                	globalTags: { env: process.env.NODE_ENV },
                	errorHandler: errorHandler,
                })
          content: |
            new StatsD($<...>)
components: []


--

