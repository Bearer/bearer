data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/breadcrumb.rb
              line_number: 2
              field_name: user_email
              object_name: metadata
    - name: Unique Identifier
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/breadcrumb.rb
              line_number: 3
              field_name: user_id
              object_name: metadata
risks:
    - detector_id: ruby_third_parties_bugsnag
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/breadcrumb.rb
              line_number: 2
              parent:
                line_number: 7
                content: Bugsnag.leave_breadcrumb('User logged in', metadata, Bugsnag::BreadcrumbType::USER)
              field_name: email
              object_name: current
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/breadcrumb.rb
              line_number: 3
              parent:
                line_number: 7
                content: Bugsnag.leave_breadcrumb('User logged in', metadata, Bugsnag::BreadcrumbType::USER)
              field_name: user_id
              object_name: metadata
components: []


--

