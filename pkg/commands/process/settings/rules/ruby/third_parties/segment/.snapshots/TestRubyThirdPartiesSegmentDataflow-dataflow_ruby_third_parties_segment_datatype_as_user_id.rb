data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/segment/testdata/datatype_as_user_id.rb
              line_number: 2
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_segment
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/segment/testdata/datatype_as_user_id.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'analytics.alias(user_id: user.email, previous_id: "some id")'
              field_name: email
              object_name: user
              subject_name: User
    - detector_id: segment_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/segment/testdata/datatype_as_user_id.rb
          line_number: 1
          parent:
            line_number: 1
            content: 'Segment::Analytics.new(write_key: "ABC123F")'
          content: |
            Segment::Analytics.new()
components: []


--

