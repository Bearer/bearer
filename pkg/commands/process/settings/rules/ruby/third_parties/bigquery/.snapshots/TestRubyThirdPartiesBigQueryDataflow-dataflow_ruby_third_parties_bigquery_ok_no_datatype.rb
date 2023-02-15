risks:
    - detector_id: ruby_third_parties_bigquery_client
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/ok_no_datatype.rb
          line_number: 1
          parent:
            line_number: 1
            content: Google::Cloud::Bigquery.new
          content: |
            Google::Cloud::Bigquery.new
    - detector_id: ruby_third_parties_bigquery_dataset
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/ok_no_datatype.rb
          line_number: 2
          parent:
            line_number: 2
            content: bigquery.dataset("my_dataset")
          content: |
            $<CLIENT>.dataset()
    - detector_id: ruby_third_parties_bigquery_insert_async
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/ok_no_datatype.rb
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

