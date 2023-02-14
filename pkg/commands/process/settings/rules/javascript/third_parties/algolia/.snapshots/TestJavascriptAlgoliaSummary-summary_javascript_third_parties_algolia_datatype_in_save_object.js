critical:
    - rule_dsrid: DSR-6
      rule_display_id: javascript_third_parties_algolia
      rule_description: Do not store sensitive data in Algolia.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_algolia
      line_number: 7
      filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_save_object.js
      category_groups:
        - PII
        - Personal Data
      parent_line_number: 8
      parent_content: |-
        index
          .saveObject(userObj, { autoGenerateObjectIDIfNotExist: true })
    - rule_dsrid: DSR-6
      rule_display_id: javascript_third_parties_algolia
      rule_description: Do not store sensitive data in Algolia.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_algolia
      line_number: 12
      filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_save_object.js
      category_groups:
        - PII
        - Personal Data
      parent_line_number: 12
      parent_content: 'index.saveObjects([{ email: user.email }])'


--

