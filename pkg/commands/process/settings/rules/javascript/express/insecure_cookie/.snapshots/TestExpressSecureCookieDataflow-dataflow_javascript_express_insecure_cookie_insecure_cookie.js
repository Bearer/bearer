risks:
    - detector_id: express_insecure_cookie
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/express/insecure_cookie/testdata/insecure_cookie.js
          line_number: 9
          parent:
            line_number: 9
            content: 'secure: false'
          content: |
            {
              cookie: {
                $<!>secure: false
              }
            }
components: []


--

