critical:
    - rule_dsrid: DSR-6
      rule_display_id: ruby_third_parties_clickhouse
      rule_description: Do not store sensitive data in ClickHouse.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_clickhouse
      line_number: 6
      filename: pkg/commands/process/settings/rules/ruby/third_parties/clickhouse/testdata/datatype_in_insert_rows.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: |-
        rows << [
              "123",
              2022,
              customer.email,
              customer.address
            ]
    - rule_dsrid: DSR-6
      rule_display_id: ruby_third_parties_clickhouse
      rule_description: Do not store sensitive data in ClickHouse.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_clickhouse
      line_number: 7
      filename: pkg/commands/process/settings/rules/ruby/third_parties/clickhouse/testdata/datatype_in_insert_rows.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: |-
        rows << [
              "123",
              2022,
              customer.email,
              customer.address
            ]


--

