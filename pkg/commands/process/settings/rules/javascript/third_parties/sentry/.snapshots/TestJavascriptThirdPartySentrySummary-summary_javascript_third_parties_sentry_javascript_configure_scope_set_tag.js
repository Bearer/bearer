critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: javascript_third_parties_sentry
      policy_description: Do not send sensitive data to Sentry.
      line_number: 2
      filename: pkg/commands/process/settings/rules/javascript/third_parties/sentry/testdata/javascript_configure_scope_set_tag.js
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: scope.setTag("user_email", user.email)


--

