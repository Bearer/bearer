data_types:
    - name: Bank Account
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_transaction_data.rb
              line_number: 1
              field_name: bank_account_number
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_analytics
      data_types:
        - name: Bank Account
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_transaction_data.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'Google::Apis::AnalyticsreportingV4::TransactionData.update!(transaction_id: "user_"+user.bank_account_number)'
              field_name: bank_account_number
              object_name: user
              subject_name: User
components: []


--

