critical:
    - policy_name: ""
      policy_dsrid: DSR-6
      policy_display_id: ruby_third_parties_elasticsearch
      policy_description: Do not store sensitive data in Elasticsearch.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/third_parties/elasticsearch/testdata/datatype_in_bulk.rb
      category_groups:
        - PII
      parent_line_number: 10
      parent_content: 'client.bulk(body: body)'


--

