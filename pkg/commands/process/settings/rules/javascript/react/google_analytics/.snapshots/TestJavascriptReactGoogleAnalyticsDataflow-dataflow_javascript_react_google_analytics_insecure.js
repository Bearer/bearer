data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/react/google_analytics/testdata/insecure.js
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/react/google_analytics/testdata/insecure.js
              line_number: 5
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: javascript_react_google_analytics
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/react/google_analytics/testdata/insecure.js
              line_number: 1
              parent:
                line_number: 2
                content: |-
                    ReactGA.event({
                    	category: "user",
                    	action: "logged_in",
                    	value: user.email,
                    })
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/react/google_analytics/testdata/insecure.js
              line_number: 5
              parent:
                line_number: 2
                content: |-
                    ReactGA.event({
                    	category: "user",
                    	action: "logged_in",
                    	value: user.email,
                    })
              field_name: email
              object_name: user
              subject_name: User
components: []


--

