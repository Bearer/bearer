data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_init.rb
              line_number: 3
              field_name: email
              object_name: current_user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_sentry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_init.rb
              line_number: 3
              parent:
                line_number: 3
                content: 'breadcrumb.message = "Authenticated user #{current_user.email}"'
              field_name: email
              object_name: current_user
              subject_name: User
components: []


--

