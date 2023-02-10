risks:
    - detector_id: javascript_http_insecure
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/lang/http_insecure/testdata/axios_insecure.js
          line_number: 2
          parent:
            line_number: 2
            content: axios.get(insecure_url)
          content: |
            $<LIBRARY>.$<METHOD>($<INSECURE_URL>)
components:
    - name: http://domain.com/api/movies
      type: ""
      sub_type: ""
      locations:
        - detector: javascript
          filename: pkg/commands/process/settings/rules/javascript/lang/http_insecure/testdata/axios_insecure.js
          line_number: 1


--

