data_types:
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/scout_apm/testdata/datatype_in_add.rb
              line_number: 1
              field_name: first_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_scout_apm
      data_types:
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/scout_apm/testdata/datatype_in_add.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'ScoutApm::Context.add({ user: { first_name: "someone" }})'
              field_name: first_name
              object_name: user
              subject_name: User
components: []


--

