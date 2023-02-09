critical:
    - rule_dsrid: DSR-1
      rule_display_id: ruby_third_parties_sentry
      rule_description: Do not send sensitive data to Sentry.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_sentry
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_breadcrumb.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        Sentry::Breadcrumb.new(
          category: "auth",
          message: "Authenticated user #{user.email}",
          level: "info"
        )


--

