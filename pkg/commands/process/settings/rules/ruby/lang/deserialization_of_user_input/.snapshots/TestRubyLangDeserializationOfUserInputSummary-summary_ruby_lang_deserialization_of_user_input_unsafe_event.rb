high:
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 2
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
      parent_line_number: 2
      parent_content: YAML.load(event["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 4
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
      parent_line_number: 4
      parent_content: Psych.load(event["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 6
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
      parent_line_number: 6
      parent_content: Syck.load(event["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 8
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
      parent_line_number: 8
      parent_content: JSON.load(event["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 10
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
      parent_line_number: 10
      parent_content: Oj.load(event["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 11
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
      parent_line_number: 11
      parent_content: |-
        Oj.object_load(event["oops"]) do |json|
          end
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 14
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
      parent_line_number: 14
      parent_content: Marshal.load(event["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_deserialization_of_user_input
      rule_description: Do not pass user input to unsafe deserialization methods.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_deserialization_of_user_input
      line_number: 15
      filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
      parent_line_number: 15
      parent_content: Marshal.restore(event["oops"])


--

