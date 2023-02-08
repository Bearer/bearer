data_types:
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_rescue_block.rb
              line_number: 4
              field_name: first_name
              object_name: current_user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_airbrake
      data_types:
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_rescue_block.rb
              line_number: 4
              parent:
                line_number: 4
                content: Airbrake.notify_sync(current_user.first_name)
              field_name: first_name
              object_name: current_user
              subject_name: User
components: []


--

