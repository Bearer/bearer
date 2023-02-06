critical:
    - policy_name: ""
      policy_dsrid: DSR-5
      policy_display_id: javascript_lang_logger
      policy_description: Do not send sensitive data to loggers.
      line_number: 1
      filename: pkg/commands/process/settings/rules/javascript/lang/logger/testdata/console.js
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: console.log(user.name);


--

