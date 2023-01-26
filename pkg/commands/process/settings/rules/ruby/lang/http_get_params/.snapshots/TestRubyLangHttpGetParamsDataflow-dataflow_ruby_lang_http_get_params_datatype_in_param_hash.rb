data_types:
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_get_params/testdata/datatype_in_param_hash.rb
              line_number: 1
              field_name: first_name
              object_name: user
risks:
    - detector_id: ruby_lang_http_get_params
      data_types:
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_get_params/testdata/datatype_in_param_hash.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'HTTP.get("https://my.api.com/users/search", params: { user: { first_name: user.first_name } })'
              field_name: first_name
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_get_params/testdata/datatype_in_param_hash.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'HTTP.get("https://my.api.com/users/search", params: { user: { first_name: user.first_name } })'
              object_name: user
components: []


--

