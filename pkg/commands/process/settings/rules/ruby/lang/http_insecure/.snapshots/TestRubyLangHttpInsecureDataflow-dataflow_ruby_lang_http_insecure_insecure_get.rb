risks:
    - detector_id: ruby_lang_http_insecure
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/insecure_get.rb
          line_number: 1
          parent:
            line_number: 1
            content: Faraday.get("http://api.insecure.com")
          content: |
            $<CLIENT>.get($<INSECURE_URL>)
components:
    - name: http://api.insecure.com
      type: ""
      sub_type: ""
      locations:
        - detector: ruby
          filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/insecure_get.rb
          line_number: 1


--

