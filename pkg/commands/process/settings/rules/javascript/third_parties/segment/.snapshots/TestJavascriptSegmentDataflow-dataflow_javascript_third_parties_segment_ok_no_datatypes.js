data_types:
    - name: Unique Identifier
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/ok_no_datatypes.js
              line_number: 4
              field_name: user_id
              object_name: track
risks:
    - detector_id: javascript_third_parties_segment_instance
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/ok_no_datatypes.js
          line_number: 2
          parent:
            line_number: 2
            content: 'new Analytics({ write_key: ''some-write-key'' })'
          content: |
            new Analytics()
components: []


--

