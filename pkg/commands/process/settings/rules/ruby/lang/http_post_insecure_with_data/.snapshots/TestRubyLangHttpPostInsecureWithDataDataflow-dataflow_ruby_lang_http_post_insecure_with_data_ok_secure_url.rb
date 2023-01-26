data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_post_insecure_with_data/testdata/ok_secure_url.rb
              line_number: 1
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/http_post_insecure_with_data/testdata/ok_secure_url.rb
              line_number: 3
              field_name: email
              object_name: user
components: []


--

