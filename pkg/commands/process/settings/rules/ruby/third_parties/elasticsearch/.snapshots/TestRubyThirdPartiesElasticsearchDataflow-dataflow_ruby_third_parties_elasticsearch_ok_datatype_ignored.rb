data_types:
    - name: Unique Identifier
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/elasticsearch/testdata/ok_datatype_ignored.rb
              line_number: 1
              field_name: user_id
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_elasticsearch
      data_types:
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/elasticsearch/testdata/ok_datatype_ignored.rb
              line_number: 1
              parent:
                line_number: 12
                content: 'client.update(index: ''books'', id: 42, body: user)'
              field_name: user_id
              object_name: user
              subject_name: User
    - detector_id: ruby_third_parties_elasticsearch_client
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/elasticsearch/testdata/ok_datatype_ignored.rb
          line_number: 2
          parent:
            line_number: 2
            content: Elasticsearch::Client.new
          content: |
            Elasticsearch::Client.new
components: []


--

