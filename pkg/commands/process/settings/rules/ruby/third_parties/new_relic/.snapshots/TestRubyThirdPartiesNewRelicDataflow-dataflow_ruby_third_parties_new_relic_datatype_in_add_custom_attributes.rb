data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_attributes.rb
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_new_relic
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_attributes.rb
              line_number: 1
              parent:
                line_number: 2
                content: NewRelic::Agent.add_custom_attributes(user)
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/datatype_in_add_custom_attributes.rb
              line_number: 3
              parent:
                line_number: 3
                content: 'NewRelic::Agent.add_custom_attributes(a: "test", user: { email: "user@example.com" }, other: 42)'
              field_name: email
              object_name: user
              subject_name: User
components: []


--

