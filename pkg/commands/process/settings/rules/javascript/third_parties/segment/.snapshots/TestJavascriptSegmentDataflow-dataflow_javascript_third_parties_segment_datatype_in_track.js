data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_track.js
              line_number: 17
              field_name: email
              object_name: user
              subject_name: User
    - name: IP address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_track.js
              line_number: 8
              field_name: ip_address
              object_name: user
              subject_name: User
    - name: Unique Identifier
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_track.js
              line_number: 7
              field_name: userId
              object_name: track
risks:
    - detector_id: javascript_third_parties_segment
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_track.js
              line_number: 17
              parent:
                line_number: 17
                content: browser.track(user.email)
              field_name: email
              object_name: user
              subject_name: User
        - name: IP address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_track.js
              line_number: 8
              parent:
                line_number: 5
                content: |-
                    client.track({
                      event: "some event name",
                      userId: user.id,
                      userIpAddr: user.ip_address,
                    })
              field_name: ip_address
              object_name: user
              subject_name: User
    - detector_id: javascript_third_parties_segment_instance
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_track.js
          line_number: 3
          parent:
            line_number: 3
            content: 'new Analytics({ write_key: ''some-write-key'' })'
          content: |
            new Analytics()
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_track.js
          line_number: 14
          parent:
            line_number: 14
            content: 'AnalyticsBrowser.load({ writeKey: ''write-key'' })'
          content: |
            AnalyticsBrowser.load()
components: []


--

