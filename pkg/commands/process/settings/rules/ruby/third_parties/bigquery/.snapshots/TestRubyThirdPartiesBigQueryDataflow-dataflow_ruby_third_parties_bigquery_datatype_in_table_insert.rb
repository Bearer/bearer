data_types:
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_table_insert.rb
              line_number: 5
              field_name: first_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_bigquery
      data_types:
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_table_insert.rb
              line_number: 5
              parent:
                line_number: 8
                content: table.insert(rows)
              field_name: first_name
              object_name: user
              subject_name: User
    - detector_id: ruby_third_parties_bigquery_client
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_table_insert.rb
          line_number: 1
          parent:
            line_number: 1
            content: Google::Cloud::Bigquery.new
          content: |
            Google::Cloud::Bigquery.new
    - detector_id: ruby_third_parties_bigquery_dataset
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_table_insert.rb
          line_number: 2
          parent:
            line_number: 2
            content: bigquery.dataset("my_dataset")
          content: |
            $<CLIENT>.dataset()
    - detector_id: ruby_third_parties_bigquery_table
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_table_insert.rb
          line_number: 3
          parent:
            line_number: 3
            content: dataset.table("my_table")
          content: |
            $<DATASET>.table()
components: []


--

