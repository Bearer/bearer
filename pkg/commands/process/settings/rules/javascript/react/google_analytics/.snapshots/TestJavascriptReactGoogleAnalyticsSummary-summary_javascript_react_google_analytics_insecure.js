critical:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_react_google_analytics
      rule_description: Do not send sensitive data to Google Analytics.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_react_google_analytics
      line_number: 1
      filename: pkg/commands/process/settings/rules/javascript/react/google_analytics/testdata/insecure.js
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: |-
        ReactGA.event({
        	category: "user",
        	action: "logged_in",
        	value: user.email,
        })
    - rule_dsrid: DSR-1
      rule_display_id: javascript_react_google_analytics
      rule_description: Do not send sensitive data to Google Analytics.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_react_google_analytics
      line_number: 5
      filename: pkg/commands/process/settings/rules/javascript/react/google_analytics/testdata/insecure.js
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: |-
        ReactGA.event({
        	category: "user",
        	action: "logged_in",
        	value: user.email,
        })


--

