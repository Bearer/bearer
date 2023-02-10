data_types:
    - name: IP address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/segment/testdata/datatype_in_nested_attribute.rb
              line_number: 2
              field_name: ip_address
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_segment
      data_types:
        - name: IP address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/segment/testdata/datatype_in_nested_attribute.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'analytics.track(user_id: user.id, event: "account sign in", context: { ip: user.ip_address })'
              field_name: ip_address
              object_name: user
              subject_name: User
    - detector_id: segment_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/segment/testdata/datatype_in_nested_attribute.rb
          line_number: 1
          parent:
            line_number: 1
            content: 'Segment::Analytics.new(write_key: "ABC123F")'
          content: |
            Segment::Analytics.new()
components: []


--

