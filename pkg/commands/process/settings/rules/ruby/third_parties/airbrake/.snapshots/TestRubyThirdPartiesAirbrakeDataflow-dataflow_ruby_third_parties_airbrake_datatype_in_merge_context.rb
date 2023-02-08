data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_merge_context.rb
              line_number: 1
              field_name: email
              object_name: current_user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_airbrake
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_merge_context.rb
              line_number: 1
              parent:
                line_number: 2
                content: 'Airbrake.merge_context(users: users)'
              field_name: email
              object_name: current_user
              subject_name: User
components: []


--

