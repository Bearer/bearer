high:
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_google_analytics
      rule_description: Do not send sensitive data to Google Analytics.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_google_analytics
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/third_parties/google_analytics/testdata/datatype_in_event_data.rb
      category_groups:
        - Personal Data
      parent_line_number: 1
      parent_content: |-
        Google::Apis::AnalyticsreportingV4::EventData.new(
          event_label: "Sign in #{customer.ip_address}"
        )


--

