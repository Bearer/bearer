critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: javascript_third_parties_sentry
      policy_description: Do not send sensitive data to Sentry.
      line_number: 1
      filename: pkg/commands/process/settings/rules/javascript/third_parties/sentry/testdata/javascript_capture_message.js
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: Sentry.captureMessage("User has successfully signed in " + current_user.email)


--

