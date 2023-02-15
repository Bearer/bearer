data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
              line_number: 9
              field_name: emailAddress
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
              line_number: 18
              field_name: email
              object_name: user
              subject_name: User
    - name: Fullname
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
              line_number: 8
              field_name: fullName
              object_name: user
              subject_name: User
    - name: Unique Identifier
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
              line_number: 6
              field_name: userId
              object_name: identify
risks:
    - detector_id: javascript_third_parties_segment
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
              line_number: 9
              parent:
                line_number: 5
                content: |-
                    analytics.identify({
                      userId: user.id,
                      traits: {
                        name: user.fullName,
                        email: user.emailAddress,
                        plan: user.businessPlan,
                        friends: user.friendCount
                      }
                    })
              field_name: emailAddress
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
              line_number: 18
              parent:
                line_number: 18
                content: browser.identify(user.email)
              field_name: email
              object_name: user
              subject_name: User
        - name: Friends
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
              line_number: 11
              parent:
                line_number: 5
                content: |-
                    analytics.identify({
                      userId: user.id,
                      traits: {
                        name: user.fullName,
                        email: user.emailAddress,
                        plan: user.businessPlan,
                        friends: user.friendCount
                      }
                    })
              field_name: friends
              object_name: traits
        - name: Fullname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
              line_number: 8
              parent:
                line_number: 5
                content: |-
                    analytics.identify({
                      userId: user.id,
                      traits: {
                        name: user.fullName,
                        email: user.emailAddress,
                        plan: user.businessPlan,
                        friends: user.friendCount
                      }
                    })
              field_name: fullName
              object_name: user
              subject_name: User
    - detector_id: javascript_third_parties_segment_instance
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
          line_number: 2
          parent:
            line_number: 2
            content: 'new Analytics({ write_key: ''some-write-key'' })'
          content: |
            new Analytics()
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
          line_number: 16
          parent:
            line_number: 16
            content: 'AnalyticsBrowser.load({ writeKey: ''write-key'' })'
          content: |
            AnalyticsBrowser.load()
components: []


--

