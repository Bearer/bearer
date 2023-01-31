risks:
    - detector_id: ruby_lang_http_insecure
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/uri_encode_form.rb
          line_number: 1
          parent:
            line_number: 1
            content: URI('http://my.api.com/users/search')
          content: |
            URI($<INSECURE_URL>)
components: []


--

