risks:
    - detector_id: express_insecure_cookie
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/express/insecure_cookie/testdata/http_only.js
          line_number: 9
          parent:
            line_number: 9
            content: 'httpOnly: false'
          content: |
            {
              cookie: {
                $<!>httpOnly: false
              }
            }
components: []


--

