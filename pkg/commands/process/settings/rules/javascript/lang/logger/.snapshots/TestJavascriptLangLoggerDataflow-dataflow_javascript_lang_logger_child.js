risks:
    - detector_id: child_logger
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/lang/logger/testdata/child.js
          line_number: 7
          parent:
            line_number: 7
            content: logger.child(ctx);
          content: |
            logger.child()
components: []


--

