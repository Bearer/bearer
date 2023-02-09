critical:
    - rule_dsrid: DSR-5
      rule_display_id: javascript_weak_encryption
      rule_description: Do not weak encrypt sensitive information
      rule_documentation_url: https://curio.sh/reference/rules/javascript_weak_encryption
      line_number: 4
      filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/md5.js
      category_groups:
        - PII
      parent_line_number: 4
      parent_content: crypto.createHmac("md5", key).update(user.password)
    - rule_dsrid: DSR-5
      rule_display_id: javascript_weak_encryption
      rule_description: Do not weak encrypt sensitive information
      rule_documentation_url: https://curio.sh/reference/rules/javascript_weak_encryption
      line_number: 5
      filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/md5.js
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: crypto.createHash("md5").update(user.password)


--

