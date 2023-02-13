data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_user_classes.rb
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
    - name: IP address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_user_classes.rb
              line_number: 4
              field_name: ip_address
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_analytics
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_user_classes.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'Google::Apis::AnalyticsreportingV4::User.new(user_id: user.email)'
              field_name: email
              object_name: user
              subject_name: User
        - name: IP address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_user_classes.rb
              line_number: 4
              parent:
                line_number: 3
                content: |-
                    Google::Apis::AnalyticsreportingV4::UserActivitySession.update!(
                      session_id: DateTime.now + user.ip_address
                    )
              field_name: ip_address
              object_name: user
              subject_name: User
components: []


--

