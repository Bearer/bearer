data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_custom_dimension.rb
              line_number: 2
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_analytics
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_custom_dimension.rb
              line_number: 2
              parent:
                line_number: 1
                content: |-
                    Google::Apis::AnalyticsreportingV4::CustomDimension.new(
                      value: user.email
                    )
              field_name: email
              object_name: user
              subject_name: User
components: []


--

