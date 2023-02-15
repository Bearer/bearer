data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/clickhouse/testdata/datatype_in_insert_rows.rb
              line_number: 6
              field_name: email
              object_name: customer
              subject_name: User
    - name: Physical Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/clickhouse/testdata/datatype_in_insert_rows.rb
              line_number: 7
              field_name: address
              object_name: customer
              subject_name: User
risks:
    - detector_id: ruby_third_parties_clickhouse
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/clickhouse/testdata/datatype_in_insert_rows.rb
              line_number: 6
              parent:
                line_number: 3
                content: |-
                    rows << [
                          "123",
                          2022,
                          customer.email,
                          customer.address
                        ]
              field_name: email
              object_name: customer
              subject_name: User
        - name: Physical Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/clickhouse/testdata/datatype_in_insert_rows.rb
              line_number: 7
              parent:
                line_number: 3
                content: |-
                    rows << [
                          "123",
                          2022,
                          customer.email,
                          customer.address
                        ]
              field_name: address
              object_name: customer
              subject_name: User
    - detector_id: ruby_third_parties_clickhouse_insert_rows
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/clickhouse/testdata/datatype_in_insert_rows.rb
          line_number: 1
          parent:
            line_number: 1
            content: rows
          content: |
            Clickhouse.connection.insert_rows() { |$<!>$<_:identifier>| }
components: []


--

