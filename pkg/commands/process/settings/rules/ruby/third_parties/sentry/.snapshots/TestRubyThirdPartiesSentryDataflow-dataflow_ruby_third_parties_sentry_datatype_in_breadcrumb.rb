data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_breadcrumb.rb
              line_number: 3
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_sentry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_breadcrumb.rb
              line_number: 3
              parent:
                line_number: 1
                content: |-
                    Sentry::Breadcrumb.new(
                      category: "auth",
                      message: "Authenticated user #{user.email}",
                      level: "info"
                    )
              field_name: email
              object_name: user
              subject_name: User
components: []


--

