critical:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_segment
      rule_description: Do not send sensitive data to Segment.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_segment
      line_number: 8
      filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: |-
        analytics.identify({
          userId: user.id,
          traits: {
            name: user.fullName,
            email: user.emailAddress,
            plan: user.businessPlan,
            friends: user.friendCount
          }
        })
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_segment
      rule_description: Do not send sensitive data to Segment.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_segment
      line_number: 9
      filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: |-
        analytics.identify({
          userId: user.id,
          traits: {
            name: user.fullName,
            email: user.emailAddress,
            plan: user.businessPlan,
            friends: user.friendCount
          }
        })
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_segment
      rule_description: Do not send sensitive data to Segment.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_segment
      line_number: 11
      filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
      category_groups:
        - PII
      parent_line_number: 5
      parent_content: |-
        analytics.identify({
          userId: user.id,
          traits: {
            name: user.fullName,
            email: user.emailAddress,
            plan: user.businessPlan,
            friends: user.friendCount
          }
        })
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_segment
      rule_description: Do not send sensitive data to Segment.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_segment
      line_number: 18
      filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_identify.js
      category_groups:
        - PII
      parent_line_number: 18
      parent_content: browser.identify(user.email)


--

