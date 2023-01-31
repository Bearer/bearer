critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_new_relic
      policy_description: Do not send sensitive data to New Relic.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_attributes.rb
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: NewRelic::Agent.add_custom_attributes(user)
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_new_relic
      policy_description: Do not send sensitive data to New Relic.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_attributes.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: 'NewRelic::Agent.add_custom_attributes(a: "test", user: { email: "user@example.com" }, other: 42)'


--

