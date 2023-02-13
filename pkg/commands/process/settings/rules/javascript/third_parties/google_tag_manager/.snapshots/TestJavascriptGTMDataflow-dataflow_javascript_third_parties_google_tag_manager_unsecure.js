data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/google_tag_manager/testdata/unsecure.js
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/google_tag_manager/testdata/unsecure.js
              line_number: 4
              field_name: email
              object_name: push
risks:
    - detector_id: javascript_google_tag_manager
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/google_tag_manager/testdata/unsecure.js
              line_number: 1
              parent:
                line_number: 3
                content: |-
                    window.dataLayer.push({
                    	email: user.email,
                    })
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/google_tag_manager/testdata/unsecure.js
              line_number: 4
              parent:
                line_number: 3
                content: |-
                    window.dataLayer.push({
                    	email: user.email,
                    })
              field_name: email
              object_name: user
              subject_name: User
components: []


--

