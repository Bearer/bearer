data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/uri_encode_form.rb
              line_number: 2
              field_name: email
              object_name: user
risks:
    - detector_id: ruby_lang_http_get
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/uri_encode_form.rb
              line_number: 2
              parent:
                line_number: 3
                content: URI.encode_www_form(user)
              field_name: email
              object_name: current_user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/uri_encode_form.rb
              line_number: 2
              parent:
                line_number: 3
                content: URI.encode_www_form(user)
              object_name: current_user
    - detector_id: ruby_lang_http_insecure
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/uri_encode_form.rb
          line_number: 1
          parent:
            line_number: 1
            content: URI('http://my.api.com/users/search')
          content: |
            URI($<INSECURE_URL>)
        - filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/uri_encode_form.rb
          line_number: 1
          parent:
            line_number: 1
            content: URI('http://my.api.com/users/search')
          content: |
            URI($<INSECURE_URL>)
components: []


--

