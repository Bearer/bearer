risks:
    - detector_id: ruby_lang_http_insecure
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/insecure_post.rb
          line_number: 1
          parent:
            line_number: 1
            content: Excon.post("http://my.api.com/users/search")
          content: |
            $<CLIENT>.post($<INSECURE_URL>)
components: []


--

