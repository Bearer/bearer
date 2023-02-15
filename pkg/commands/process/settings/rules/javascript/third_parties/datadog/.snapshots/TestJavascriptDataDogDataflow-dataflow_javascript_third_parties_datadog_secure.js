risks:
    - detector_id: javascript_third_parties_hotshot_statsd
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/datadog/testdata/secure.js
          line_number: 2
          parent:
            line_number: 2
            content: |-
                new StatsD({
                	port: 8020,
                	globalTags: { env: process.env.NODE_ENV },
                	errorHandler: errorHandler,
                })
          content: |
            new StatsD($<...>)
components: []


--

