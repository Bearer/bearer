low:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_rollbar
      policy_description: Do not send sensitive data to Rollbar.
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/third_parties/rollbar/testdata/datatype_in_scoped.rb
      parent_line_number: 3
      parent_content: |-
        Rollbar.scoped(scope) do
          call
        end


--

