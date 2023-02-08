critical:
    - policy_name: ""
      policy_dsrid: DSR-1
      policy_display_id: ruby_third_parties_bugsnag
      policy_description: Do not send sensitive data to Bugsnag.
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/breadcrumb.rb
      category_groups:
        - PII
      parent_line_number: 7
      parent_content: Bugsnag.leave_breadcrumb('User logged in', metadata, Bugsnag::BreadcrumbType::USER)


--

