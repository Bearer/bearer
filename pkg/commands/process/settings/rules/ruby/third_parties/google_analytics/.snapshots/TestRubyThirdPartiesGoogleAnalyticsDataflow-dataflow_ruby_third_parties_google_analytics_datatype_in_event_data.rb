data_types:
    - name: IP address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_event_data.rb
              line_number: 2
              field_name: ip_address
              object_name: customer
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_analytics
      data_types:
        - name: IP address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_event_data.rb
              line_number: 2
              parent:
                line_number: 1
                content: |-
                    Google::Apis::AnalyticsreportingV4::EventData.new(
                      event_label: "Sign in #{customer.ip_address}"
                    )
              field_name: ip_address
              object_name: customer
              subject_name: User
components: []


--

