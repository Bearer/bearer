data_types:
    - name: Unique Identifier
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/ok_datatype_ignored.rb
              line_number: 3
              field_name: user_id
              object_name: current_user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_sentry
      data_types:
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/ok_datatype_ignored.rb
              line_number: 3
              parent:
                line_number: 1
                content: |-
                    Sentry::Breadcrumb.new(
                      category: "auth",
                      message: "user has authenticated #{current_user.user_id}",
                      level: "info"
                    )
              field_name: user_id
              object_name: current_user
              subject_name: User
components: []


--

