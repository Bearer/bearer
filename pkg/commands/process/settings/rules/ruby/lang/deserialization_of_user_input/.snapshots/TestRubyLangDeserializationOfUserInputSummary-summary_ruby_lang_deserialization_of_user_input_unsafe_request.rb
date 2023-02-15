high:
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
      parent_line_number: 1
      parent_content: YAML.load(request.env[:oops])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
      parent_line_number: 3
      parent_content: Psych.load(request.env[:oops])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 5
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
      parent_line_number: 5
      parent_content: Syck.load(request.env[:oops])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 7
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
      parent_line_number: 7
      parent_content: JSON.load(request.env[:oops])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 9
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
      parent_line_number: 9
      parent_content: Oj.load(request.env[:oops])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 10
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
      parent_line_number: 10
      parent_content: |-
        Oj.object_load(request.env[:oops]) do |json|
        end
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 13
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
      parent_line_number: 13
      parent_content: Marshal.load(request.env[:oops])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 14
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
      parent_line_number: 14
      parent_content: Marshal.restore(request.env[:oops])


--

