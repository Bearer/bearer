data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/before_bugsnag_notify_add_metadata.rb
              line_number: 5
              field_name: email
              object_name: user
              subject_name: User
    - name: Fullname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/before_bugsnag_notify_add_metadata.rb
              line_number: 12
              field_name: name
              object_name: current_user
              subject_name: User
components: []


--

