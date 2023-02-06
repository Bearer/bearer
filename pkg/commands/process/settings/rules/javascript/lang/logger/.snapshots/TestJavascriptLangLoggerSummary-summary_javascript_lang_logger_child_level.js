critical:
    - policy_name: ""
      policy_dsrid: DSR-5
      policy_display_id: javascript_lang_logger
      policy_description: Do not send sensitive data to loggers.
      line_number: 3
      filename: pkg/commands/process/settings/rules/javascript/lang/logger/testdata/child_level.js
      category_groups:
        - PII
      parent_line_number: 7
      parent_content: logger.child(ctx)
    - policy_name: ""
      policy_dsrid: DSR-5
      policy_display_id: javascript_lang_logger
      policy_description: Do not send sensitive data to loggers.
      line_number: 7
      filename: pkg/commands/process/settings/rules/javascript/lang/logger/testdata/child_level.js
      category_groups:
        - PII
      parent_line_number: 7
      parent_content: logger.child(ctx).info(user.name)


--

