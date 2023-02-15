data_types:
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_insert.rb
              line_number: 4
              field_name: first_name
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_bigquery
      data_types:
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_insert.rb
              line_number: 4
              parent:
                line_number: 7
                content: |-
                    dataset.insert("my_table", rows) do |result|
                      call
                    end
              field_name: first_name
              object_name: user
              subject_name: User
    - detector_id: ruby_third_parties_bigquery_client
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_insert.rb
          line_number: 1
          parent:
            line_number: 1
            content: 'Google::Cloud::Bigquery.new(retries: 3)'
          content: |
            Google::Cloud::Bigquery.new()
    - detector_id: ruby_third_parties_bigquery_dataset
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bigquery/testdata/datatype_in_insert.rb
          line_number: 2
          parent:
            line_number: 2
            content: bigquery.dataset("my_dataset")
          content: |
            $<CLIENT>.dataset()
components: []


--

