data_types:
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_notice_error.rb
              line_number: 1
              field_name: first_name
              object_name: user
risks:
    - detector_id: ruby_third_parties_new_relic
      data_types:
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_notice_error.rb
              line_number: 1
              parent:
                line_number: 2
                content: 'NewRelic::Agent.notice_error(exception, { custom_params: user })'
              field_name: first_name
              object_name: user
        - name: Lastname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_notice_error.rb
              line_number: 3
              parent:
                line_number: 3
                content: 'NewRelic::Agent.notice_error(exception, expected: true, custom_params: { last_name: "foo" }, metric: "test")'
              field_name: last_name
              object_name: custom_params
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_notice_error.rb
              line_number: 1
              parent:
                line_number: 2
                content: 'NewRelic::Agent.notice_error(exception, { custom_params: user })'
              object_name: user
components: []


--

