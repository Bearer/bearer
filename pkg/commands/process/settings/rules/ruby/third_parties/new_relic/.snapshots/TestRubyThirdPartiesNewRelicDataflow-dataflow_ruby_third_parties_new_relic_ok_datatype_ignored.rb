data_types:
    - name: Unique Identifier
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/ok_datatype_ignored.rb
              line_number: 1
              field_name: user_id
              object_name: user
risks:
    - detector_id: ruby_third_parties_new_relic
      data_types:
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/new_relic/testdata/ok_datatype_ignored.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'NewRelic::Agent.add_custom_attributes(user_id: user.user_id)'
              field_name: user_id
              object_name: user
components: []


--

