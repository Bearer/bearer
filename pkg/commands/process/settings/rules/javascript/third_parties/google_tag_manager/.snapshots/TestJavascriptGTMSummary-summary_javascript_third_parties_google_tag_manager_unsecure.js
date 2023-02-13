critical:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_google_tag_manager
      rule_description: Do not send sensitive data to google tag manager.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_google_tag_manager
      line_number: 1
      filename: pkg/commands/process/settings/rules/javascript/third_parties/google_tag_manager/testdata/unsecure.js
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: |-
        window.dataLayer.push({
        	email: user.email,
        })
    - rule_dsrid: DSR-1
      rule_display_id: javascript_google_tag_manager
      rule_description: Do not send sensitive data to google tag manager.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_google_tag_manager
      line_number: 4
      filename: pkg/commands/process/settings/rules/javascript/third_parties/google_tag_manager/testdata/unsecure.js
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: |-
        window.dataLayer.push({
        	email: user.email,
        })


--

