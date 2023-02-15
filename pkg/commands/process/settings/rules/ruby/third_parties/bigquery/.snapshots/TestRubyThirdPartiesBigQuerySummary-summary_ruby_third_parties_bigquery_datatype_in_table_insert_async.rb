critical:
    - rule_dsrid: DSR-6
      rule_display_id: ruby_third_parties_bigquery
      rule_description: Do not store sensitive data in BigQuery.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_bigquery
      line_number: 9
      filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_table_insert_async.rb
      category_groups:
        - PII
      parent_line_number: 11
      parent_content: inserter.insert(rows)


--

