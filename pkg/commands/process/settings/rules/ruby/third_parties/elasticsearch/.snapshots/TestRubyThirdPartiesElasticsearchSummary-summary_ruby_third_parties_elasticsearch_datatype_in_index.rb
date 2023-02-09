critical:
    - rule_dsrid: DSR-6
      rule_display_id: ruby_third_parties_elasticsearch
      rule_description: Do not store sensitive data in Elasticsearch.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_elasticsearch
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/third_parties/elasticsearch/testdata/datatype_in_index.rb
      category_groups:
        - PII
      parent_line_number: 4
      parent_content: 'client.index({ index: ''users'', body: user })'


--

