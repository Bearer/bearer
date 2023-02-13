low:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_google_analytics
      rule_description: Do not send sensitive data to Google Analytics.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_google_analytics
      line_number: 3
      filename: pkg/commands/process/settings/rules/javascript/third_parties/google_analytics/testdata/unsecure.js
      parent_line_number: 1
      parent_content: |-
        gtag("event", "screen_view", {
        	user: {
        		email: "jhon@gmail.com",
        	},
        })


--

