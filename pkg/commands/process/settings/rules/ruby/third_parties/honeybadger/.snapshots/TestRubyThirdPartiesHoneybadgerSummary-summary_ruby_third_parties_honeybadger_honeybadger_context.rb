critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_honeybadger
      policy_description: Do not send sensitive data to Honeybadger.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_context.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: |-
        Honeybadger.context({
          tags: tags
        })
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_honeybadger
      policy_description: Do not send sensitive data to Honeybadger.
      line_number: 8
      filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_context.rb
      category_groups:
        - PII
      parent_line_number: 7
      parent_content: |-
        Honeybadger.context({
          my_data: current_user.email
        })


--

