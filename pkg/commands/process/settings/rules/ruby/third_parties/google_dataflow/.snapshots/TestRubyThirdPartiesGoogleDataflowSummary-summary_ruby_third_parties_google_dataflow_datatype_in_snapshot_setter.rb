critical:
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_google_dataflow
      rule_description: Do not send sensitive data to Google Dataflow.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_google_dataflow
      line_number: 9
      filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_snapshot_setter.rb
      category_groups:
        - PII
      parent_line_number: 9
      parent_content: 'snapshot.description = "Current user: #{user.email_address}"'


--

