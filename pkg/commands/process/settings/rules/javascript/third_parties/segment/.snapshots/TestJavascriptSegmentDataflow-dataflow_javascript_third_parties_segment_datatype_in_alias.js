data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_alias.js
              line_number: 8
              field_name: email
              object_name: user
              subject_name: User
    - name: Unique Identifier
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_alias.js
              line_number: 9
              field_name: userId
              object_name: alias
risks:
    - detector_id: javascript_third_parties_segment
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_alias.js
              line_number: 8
              parent:
                line_number: 7
                content: |-
                    appAnalytics.alias({
                      previousId: user.email,
                      userId: user.id,
                    })
              field_name: email
              object_name: user
              subject_name: User
    - detector_id: javascript_third_parties_segment_instance
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_alias.js
          line_number: 3
          parent:
            line_number: 3
            content: 'new Analytics({ writeKey: "product-write-key" })'
          content: |
            new Analytics()
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_alias.js
          line_number: 4
          parent:
            line_number: 4
            content: 'new Analytics({ writeKey: "application-write-key" })'
          content: |
            new Analytics()
components: []


--

