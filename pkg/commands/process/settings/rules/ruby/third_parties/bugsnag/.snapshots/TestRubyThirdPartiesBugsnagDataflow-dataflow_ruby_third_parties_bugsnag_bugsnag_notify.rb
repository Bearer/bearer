data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/bugsnag_notify.rb
              line_number: 2
              field_name: email
              object_name: current_user
              subject_name: User
    - name: Fullname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/bugsnag_notify.rb
              line_number: 13
              field_name: name
              object_name: current_user
              subject_name: User
components: []


--

