critical:
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_new_relic
      rule_description: Do not send sensitive data to New Relic.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_new_relic
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_attributes.rb
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: NewRelic::Agent.add_custom_attributes(user)
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_new_relic
      rule_description: Do not send sensitive data to New Relic.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_new_relic
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_attributes.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: 'NewRelic::Agent.add_custom_attributes(a: "test", user: { email: "user@example.com" }, other: 42)'


--

