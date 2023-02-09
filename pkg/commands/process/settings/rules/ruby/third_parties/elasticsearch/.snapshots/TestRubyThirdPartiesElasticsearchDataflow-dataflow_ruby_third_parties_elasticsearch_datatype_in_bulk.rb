data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/elasticsearch/testdata/datatype_in_bulk.rb
              line_number: 3
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_elasticsearch
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/elasticsearch/testdata/datatype_in_bulk.rb
              line_number: 3
              parent:
                line_number: 10
                content: 'client.bulk(body: body)'
              field_name: email
              object_name: user
              subject_name: User
    - detector_id: ruby_third_parties_elasticsearch_client
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/elasticsearch/testdata/datatype_in_bulk.rb
          line_number: 1
          parent:
            line_number: 1
            content: 'Elasticsearch::Client.new(log: true)'
          content: |
            Elasticsearch::Client.new()
components: []


--

