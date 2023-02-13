data_types:
    - name: Physical Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_cohort.rb
              line_number: 1
              field_name: zip_code
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_analytics
      data_types:
        - name: Physical Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_cohort.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'Google::Apis::AnalyticsreportingV4::Cohort.new(name: user.zip_code)'
              field_name: zip_code
              object_name: user
              subject_name: User
components: []


--

