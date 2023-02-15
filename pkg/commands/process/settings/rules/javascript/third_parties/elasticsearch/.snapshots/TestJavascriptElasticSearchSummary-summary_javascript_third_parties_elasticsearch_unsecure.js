critical:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_elasticsearch
      rule_description: Do not send sensitive data to ElasticSearch.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_elasticsearch
      line_number: 1
      filename: pkg/commands/process/settings/rules/javascript/third_parties/elasticsearch/testdata/unsecure.js
      category_groups:
        - PII
      parent_line_number: 2
      parent_content: elasticsearch.index(user)


--

