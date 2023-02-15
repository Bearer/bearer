data_types:
    - name: Unique Identifier
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/ok_ignored_datatypes.js
              line_number: 7
              field_name: userId
              object_name: track
risks:
    - detector_id: javascript_third_parties_segment_instance
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/ok_ignored_datatypes.js
          line_number: 3
          parent:
            line_number: 3
            content: 'new Analytics({ write_key: ''some-write-key'' })'
          content: |
            new Analytics()
components: []


--

