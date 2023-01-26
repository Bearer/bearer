critical:
    - policy_name: ""
      policy_dsrid: DSR-5
      policy_display_id: ruby_rails_logger
      policy_description: Do not send sensitive data to loggers.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/rails/logger/testdata/datatype_leak.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: Rails.logger.info(user.email)


--

