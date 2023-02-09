risks:
    - detector_id: segment_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/segment/testdata/ok_only_ignored_datatypes.rb
          line_number: 1
          parent:
            line_number: 1
            content: 'Segment::Analytics.new(write_key: "ABC123F")'
          content: |
            Segment::Analytics.new()
components: []


--

