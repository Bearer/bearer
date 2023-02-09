risks:
    - detector_id: ruby_third_parties_elasticsearch_client
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/elasticsearch/testdata/ok_no_datatype.rb
          line_number: 2
          parent:
            line_number: 2
            content: Elasticsearch::Client.new
          content: |
            Elasticsearch::Client.new
components: []


--

