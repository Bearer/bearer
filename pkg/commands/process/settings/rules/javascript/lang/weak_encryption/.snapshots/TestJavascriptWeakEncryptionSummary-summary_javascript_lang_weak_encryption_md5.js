critical:
    - policy_name: ""
      policy_dsrid: DSR-5
      policy_display_id: javascript_weak_encryption
      policy_description: Do not weak encrypt sensitive information
      line_number: 4
      filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/md5.js
      category_groups:
        - PII
      parent_line_number: 4
      parent_content: crypto.createHmac("md5", key).update(user.password)
    - policy_name: ""
      policy_dsrid: DSR-5
      policy_display_id: javascript_weak_encryption
      policy_description: Do not weak encrypt sensitive information
      line_number: 5
      filename: pkg/commands/process/settings/rules/javascript/lang/weak_encryption/testdata/md5.js
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: crypto.createHash("md5").update(user.password)


--

