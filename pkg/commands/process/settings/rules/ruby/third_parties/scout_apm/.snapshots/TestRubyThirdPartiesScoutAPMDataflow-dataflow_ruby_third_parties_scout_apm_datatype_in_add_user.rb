data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/scout_apm/testdata/datatype_in_add_user.rb
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_scout_apm
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/scout_apm/testdata/datatype_in_add_user.rb
              line_number: 1
              parent:
                line_number: 2
                content: ScoutApm::Context.add_user(user)
              field_name: email
              object_name: user
              subject_name: User
components: []


--

