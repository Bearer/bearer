critical:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_new_relic
      rule_description: Do not send sensitive data to New Relic.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_new_relic
      line_number: 6
      filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_interaction_set_attribute.js
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: |-
        newrelic.interaction()
            .setAttribute("username", user.first_name)
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_new_relic
      rule_description: Do not send sensitive data to New Relic.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_new_relic
      line_number: 13
      filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_interaction_set_attribute.js
      category_groups:
        - PII
      parent_line_number: 13
      parent_content: interaction.setAttribute("email", user.email_address)


--

