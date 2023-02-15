risks:
    - detector_id: ruby_third_parties_clickhouse_insert_rows
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/clickhouse/testdata/ok_no_datatype.rb
          line_number: 1
          parent:
            line_number: 1
            content: rows
          content: |
            Clickhouse.connection.insert_rows() { |$<!>$<_:identifier>| }
components: []


--

