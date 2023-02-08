data_types:
    - name: Unique Identifier
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/scout_apm/testdata/ok_datatype_ignored.rb
              line_number: 1
              field_name: user_id
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_scout_apm
      data_types:
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/scout_apm/testdata/ok_datatype_ignored.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'ScoutApm::Context.add({ user: { user_id: 42  } })'
              field_name: user_id
              object_name: user
              subject_name: User
components: []


--

