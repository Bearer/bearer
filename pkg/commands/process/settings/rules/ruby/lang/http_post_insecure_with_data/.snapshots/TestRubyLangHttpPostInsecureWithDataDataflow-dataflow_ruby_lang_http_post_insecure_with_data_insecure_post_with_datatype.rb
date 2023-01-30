data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_post_insecure_with_data/testdata/insecure_post_with_datatype.rb
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_lang_http_post_insecure_with_data
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_post_insecure_with_data/testdata/insecure_post_with_datatype.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'HTTParty.post("http://my.api.com/users/search", body: user.email)'
              field_name: email
              object_name: user
              subject_name: User
    - detector_id: ruby_lang_http_insecure
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/http_post_insecure_with_data/testdata/insecure_post_with_datatype.rb
          line_number: 1
          parent:
            line_number: 1
            content: 'HTTParty.post("http://my.api.com/users/search", body: user.email)'
          content: |
            $<CLIENT>.post($<INSECURE_URL>)
components: []


--

