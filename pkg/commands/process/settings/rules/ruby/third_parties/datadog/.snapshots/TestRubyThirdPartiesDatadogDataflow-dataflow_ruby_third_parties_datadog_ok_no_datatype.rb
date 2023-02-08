risks:
    - detector_id: ruby_third_parties_datadog_span
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/datadog/testdata/ok_no_datatype.rb
          line_number: 6
          parent:
            line_number: 6
            content: span
          content: |
            Datadog::Tracing.trace() { |$<!>$<SPAN:identifier>$<...>| }
components: []


--

