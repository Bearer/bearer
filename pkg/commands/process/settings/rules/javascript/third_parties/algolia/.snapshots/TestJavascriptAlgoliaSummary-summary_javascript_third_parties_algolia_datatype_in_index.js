critical:
    - rule_dsrid: DSR-6
      rule_display_id: javascript_third_parties_algolia
      rule_description: Do not store sensitive data in Algolia.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_algolia
      line_number: 4
      filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_index.js
      category_groups:
        - PII
      parent_line_number: 4
      parent_content: myAlgolia.initIndex(company.bank_account_number)


--

