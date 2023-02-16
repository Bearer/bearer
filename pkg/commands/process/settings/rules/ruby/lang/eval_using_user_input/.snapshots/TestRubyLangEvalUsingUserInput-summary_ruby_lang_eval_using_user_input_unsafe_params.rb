high:
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_eval_using_user_input
      rule_description: Do not generate code using user input.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_eval_using_user_input
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/eval_using_user_input/testdata/unsafe_params.rb
      parent_line_number: 1
      parent_content: RubyVM::InstructionSequence.compile(params["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_eval_using_user_input
      rule_description: Do not generate code using user input.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_eval_using_user_input
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/lang/eval_using_user_input/testdata/unsafe_params.rb
      parent_line_number: 3
      parent_content: a.eval(params["oops"], "test")
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_eval_using_user_input
      rule_description: Do not generate code using user input.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_eval_using_user_input
      line_number: 5
      filename: pkg/commands/process/settings/rules/ruby/lang/eval_using_user_input/testdata/unsafe_params.rb
      parent_line_number: 5
      parent_content: a.instance_eval(params["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_eval_using_user_input
      rule_description: Do not generate code using user input.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_eval_using_user_input
      line_number: 7
      filename: pkg/commands/process/settings/rules/ruby/lang/eval_using_user_input/testdata/unsafe_params.rb
      parent_line_number: 7
      parent_content: a.class_eval(params["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_eval_using_user_input
      rule_description: Do not generate code using user input.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_eval_using_user_input
      line_number: 9
      filename: pkg/commands/process/settings/rules/ruby/lang/eval_using_user_input/testdata/unsafe_params.rb
      parent_line_number: 9
      parent_content: a.module_eval(params["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_eval_using_user_input
      rule_description: Do not generate code using user input.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_eval_using_user_input
      line_number: 11
      filename: pkg/commands/process/settings/rules/ruby/lang/eval_using_user_input/testdata/unsafe_params.rb
      parent_line_number: 11
      parent_content: eval(params["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_eval_using_user_input
      rule_description: Do not generate code using user input.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_eval_using_user_input
      line_number: 13
      filename: pkg/commands/process/settings/rules/ruby/lang/eval_using_user_input/testdata/unsafe_params.rb
      parent_line_number: 13
      parent_content: instance_eval(params["oops"], "test")
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_eval_using_user_input
      rule_description: Do not generate code using user input.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_eval_using_user_input
      line_number: 15
      filename: pkg/commands/process/settings/rules/ruby/lang/eval_using_user_input/testdata/unsafe_params.rb
      parent_line_number: 15
      parent_content: class_eval(params["oops"])
    - rule_dsrid: DSR-?
      rule_display_id: ruby_lang_eval_using_user_input
      rule_description: Do not generate code using user input.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_eval_using_user_input
      line_number: 17
      filename: pkg/commands/process/settings/rules/ruby/lang/eval_using_user_input/testdata/unsafe_params.rb
      parent_line_number: 17
      parent_content: module_eval(params["oops"])


--

