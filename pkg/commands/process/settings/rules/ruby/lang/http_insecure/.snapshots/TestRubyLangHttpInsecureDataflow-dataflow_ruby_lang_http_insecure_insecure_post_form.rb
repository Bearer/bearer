risks:
    - detector_id: ruby_lang_http_insecure
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/insecure_post_form.rb
          line_number: 1
          parent:
            line_number: 1
            content: Net::HTTP.post_form("http://my.api.com/users/search")
          content: |
            Net::HTTP.post_form($<INSECURE_URL>)
        - filename: pkg/commands/process/settings/rules/ruby/lang/http_insecure/testdata/insecure_post_form.rb
          line_number: 1
          parent:
            line_number: 1
            content: Net::HTTP.post_form("http://my.api.com/users/search")
          content: |
            Net::HTTP.post_form($<INSECURE_URL>)
components: []


--

