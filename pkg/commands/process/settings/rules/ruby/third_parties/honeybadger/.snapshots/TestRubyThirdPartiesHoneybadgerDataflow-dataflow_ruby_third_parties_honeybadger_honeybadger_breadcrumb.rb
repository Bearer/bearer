data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_breadcrumb.rb
              line_number: 1
              field_name: email
              object_name: current_user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_honeybadger
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/honeybadger/testdata/honeybadger_breadcrumb.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'Honeybadger.add_breadcrumb("Email Sent", metadata: { user_email: current_user.email, message: message })'
              field_name: email
              object_name: current_user
              subject_name: User
components: []


--

