risks:
    - detector_id: javascript_http_insecure
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/lang/http_insecure/testdata/request_insecure.js
          line_number: 5
          parent:
            line_number: 5
            content: xhttp.open("GET", insecure_url, true)
          content: |
            $<REQUEST>.open($<_>, $<INSECURE_URL>);
components: []


--

