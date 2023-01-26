data_types:
    - name: Ethnic Origin
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_get_params/testdata/datatype_in_params.rb
              line_number: 1
              field_name: ethnic_origin
              object_name: user
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_get_params/testdata/datatype_in_params.rb
              line_number: 3
              field_name: first_name
              object_name: user
risks:
    - detector_id: ruby_lang_http_get_params
      data_types:
        - name: Ethnic Origin
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_get_params/testdata/datatype_in_params.rb
              line_number: 1
              parent:
                line_number: 1
                content: URI("https://my.api.com/users/search?ethnic_origin=#{user.ethnic_origin}")
              field_name: ethnic_origin
              object_name: user
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_get_params/testdata/datatype_in_params.rb
              line_number: 3
              parent:
                line_number: 3
                content: RestClient.get("https://my.api.com/users/search?first_name=#{user.first_name}")
              field_name: first_name
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_get_params/testdata/datatype_in_params.rb
              line_number: 1
              parent:
                line_number: 1
                content: URI("https://my.api.com/users/search?ethnic_origin=#{user.ethnic_origin}")
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_get_params/testdata/datatype_in_params.rb
              line_number: 3
              parent:
                line_number: 3
                content: RestClient.get("https://my.api.com/users/search?first_name=#{user.first_name}")
              object_name: user
components: []


--

