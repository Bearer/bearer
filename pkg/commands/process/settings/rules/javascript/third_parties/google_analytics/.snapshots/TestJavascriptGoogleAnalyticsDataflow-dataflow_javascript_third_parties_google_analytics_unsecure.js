risks:
    - detector_id: javascript_google_analytics
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/google_analytics/testdata/unsecure.js
              line_number: 3
              parent:
                line_number: 1
                content: |-
                    gtag("event", "screen_view", {
                    	user: {
                    		email: "jhon@gmail.com",
                    	},
                    })
              field_name: email
              object_name: user
              subject_name: User
components: []


--

