critical:
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_google_analytics
      rule_description: Do not send sensitive data to Google Analytics.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_google_analytics
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_user_classes.rb
      category_groups:
        - PII
        - Personal Data
      parent_line_number: 1
      parent_content: 'Google::Apis::AnalyticsreportingV4::User.new(user_id: user.email)'
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_google_analytics
      rule_description: Do not send sensitive data to Google Analytics.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_google_analytics
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_user_classes.rb
      category_groups:
        - PII
        - Personal Data
      parent_line_number: 3
      parent_content: |-
        Google::Apis::AnalyticsreportingV4::UserActivitySession.update!(
          session_id: DateTime.now + user.ip_address
        )


--

