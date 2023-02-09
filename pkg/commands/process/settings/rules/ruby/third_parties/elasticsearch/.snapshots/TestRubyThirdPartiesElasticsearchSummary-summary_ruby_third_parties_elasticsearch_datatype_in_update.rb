critical:
    - rule_dsrid: DSR-6
      rule_display_id: ruby_third_parties_elasticsearch
      rule_description: Do not store sensitive data in Elasticsearch.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_third_parties_elasticsearch
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/third_parties/elasticsearch/testdata/datatype_in_update.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: |-
        Elasticsearch::Client
          .new
          .update(index: 'books', id: 42, body: user)


--

