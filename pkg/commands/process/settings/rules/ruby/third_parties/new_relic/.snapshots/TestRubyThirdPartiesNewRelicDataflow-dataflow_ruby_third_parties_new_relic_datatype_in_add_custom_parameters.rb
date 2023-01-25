data_types:
    - name: Physical Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_parameters.rb
              line_number: 1
              field_name: address
              object_name: user
risks:
    - detector_id: ruby_third_parties_new_relic
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_parameters.rb
              line_number: 3
              parent:
                line_number: 3
                content: 'NewRelic::Agent.add_custom_parameters(user: { email: "user@example.com" })'
              field_name: email
              object_name: user
        - name: Physical Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_parameters.rb
              line_number: 1
              parent:
                line_number: 2
                content: NewRelic::Agent.add_custom_parameters(user)
              field_name: address
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_parameters.rb
              line_number: 1
              parent:
                line_number: 2
                content: NewRelic::Agent.add_custom_parameters(user)
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_parameters.rb
              line_number: 3
              parent:
                line_number: 3
                content: 'NewRelic::Agent.add_custom_parameters(user: { email: "user@example.com" })'
              object_name: user
components: []


--

