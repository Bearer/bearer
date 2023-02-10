critical:
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_segment
      rule_description: Do not send sensitive data to Segment.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_segment
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/third_parties/segment/testdata/datatype_as_user_id.rb
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: 'analytics.alias(user_id: user.email, previous_id: "some id")'


--

