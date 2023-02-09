critical:
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_sentry
      rule_description: Do not send sensitive data to Sentry.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_sentry
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_tag.rb
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: scope.set_tag(:email, user.email)
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_sentry
      rule_description: Do not send sensitive data to Sentry.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_sentry
      line_number: 6
      filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_set_tag.rb
      category_groups:
        - PII
      parent_line_number: 6
      parent_content: scope.set_tag(:email, user.email)


--

