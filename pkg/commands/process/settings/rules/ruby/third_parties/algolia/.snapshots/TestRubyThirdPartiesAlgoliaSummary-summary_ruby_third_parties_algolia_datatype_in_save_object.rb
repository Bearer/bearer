critical:
    - policy_name: ""
      policy_dsrid: DSR-6
      policy_display_id: ruby_third_parties_algolia
      policy_description: Do not store sensitive data in Algolia.
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/third_parties/algolia/testdata/datatype_in_save_object.rb
      category_groups:
        - PII
      parent_line_number: 4
      parent_content: 'index.save_object({ email: user.email }, { auto_generate_object_id_if_not_exist: true })'
    - policy_name: ""
      policy_dsrid: DSR-6
      policy_display_id: ruby_third_parties_algolia
      policy_description: Do not store sensitive data in Algolia.
      line_number: 6
      filename: pkg/commands/process/settings/rules/ruby/third_parties/algolia/testdata/datatype_in_save_object.rb
      category_groups:
        - PII
      parent_line_number: 6
      parent_content: 'index.save_objects([{ email: user.email }], { auto_generate_object_id_if_not_exist: true })'


--

