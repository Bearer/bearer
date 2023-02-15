risks:
    - detector_id: javascript_third_parties_datadog_browser
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/datadog_browser/testdata/unsecure.js
          line_number: 2
          parent:
            line_number: 2
            content: 'trackUserInteractions: true'
          content: |
            DD_RUM.init({
                $<!>trackUserInteractions: true,
              })
components: []


--

