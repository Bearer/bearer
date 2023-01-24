critical:
    - policy_name: ""
      policy_dsrid: DSR-5
      policy_display_id: ruby_lang_logger
      policy_description: Do not send sensitive data to loggers.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/lang/logger/testdata/datatype_leak.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        logger.info(
          "user info:",
          user.email,
          user.address
        )
    - policy_name: ""
      policy_dsrid: DSR-5
      policy_display_id: ruby_lang_logger
      policy_description: Do not send sensitive data to loggers.
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/lang/logger/testdata/datatype_leak.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        logger.info(
          "user info:",
          user.email,
          user.address
        )


--

