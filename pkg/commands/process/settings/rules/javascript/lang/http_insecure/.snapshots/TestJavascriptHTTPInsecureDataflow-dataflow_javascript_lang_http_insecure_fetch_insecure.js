risks:
    - detector_id: javascript_http_insecure
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/lang/http_insecure/testdata/fetch_insecure.js
          line_number: 3
          parent:
            line_number: 3
            content: fetch(insecure_url)
          content: |
            fetch($<INSECURE_URL>)
components: []


--

