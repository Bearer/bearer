data_types:
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_insert_async.rb
              line_number: 8
              field_name: :first_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_bigquery
      data_types:
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_insert_async.rb
              line_number: 8
              parent:
                line_number: 10
                content: inserter.insert(rows)
              field_name: first_name
              object_name: user
              subject_name: User
    - detector_id: ruby_third_parties_bigquery_client
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_insert_async.rb
          line_number: 1
          parent:
            line_number: 1
            content: Google::Cloud::Bigquery.new
          content: |
            Google::Cloud::Bigquery.new
    - detector_id: ruby_third_parties_bigquery_dataset
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_insert_async.rb
          line_number: 2
          parent:
            line_number: 2
            content: bigquery.dataset("my_dataset")
          content: |
            $<CLIENT>.dataset()
    - detector_id: ruby_third_parties_bigquery_insert_async
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_insert_async.rb
          line_number: 4
          parent:
            line_number: 4
            content: |-
                dataset.insert_async "my_table" do |result|
                  call
                end
          content: |
            $<DATASET>.insert_async()$<...>
components: []


--

