data_types:
    - name: Job Titles
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_group.js
              line_number: 8
              field_name: job_title
              object_name: user
              subject_name: User
    - name: Unique Identifier
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_group.js
              line_number: 7
              field_name: userId
              object_name: group
risks:
    - detector_id: javascript_third_parties_segment
      data_types:
        - name: Job Titles
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_group.js
              line_number: 8
              parent:
                line_number: 6
                content: |-
                    analytics.group({
                      userId: user.id,
                      groupId: user.job_title,
                      traits: {},
                    })
              field_name: job_title
              object_name: user
              subject_name: User
    - detector_id: javascript_third_parties_segment_instance
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_group.js
          line_number: 3
          parent:
            line_number: 3
            content: 'new Analytics({ writeKey: ''my-write-key'' })'
          content: |
            new Analytics()
components: []


--

