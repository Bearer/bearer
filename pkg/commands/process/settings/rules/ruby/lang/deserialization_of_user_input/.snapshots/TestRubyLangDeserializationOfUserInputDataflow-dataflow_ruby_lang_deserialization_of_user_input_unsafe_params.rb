risks:
    - detector_id: ruby_lang_deserialization_of_user_input
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 1
          parent:
            line_number: 1
            content: YAML.load(params[:oops])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 3
          parent:
            line_number: 3
            content: Psych.load(params[:oops])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 5
          parent:
            line_number: 5
            content: Syck.load(params[:oops])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 7
          parent:
            line_number: 7
            content: JSON.load(params[:oops])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 9
          parent:
            line_number: 9
            content: |-
                Oj.load(params[:oops]) do |json|
                end
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 11
          parent:
            line_number: 11
            content: Oj.object_load(params[:oops])
          content: |
            Oj.object_load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 13
          parent:
            line_number: 13
            content: Marshal.load(params[:oops])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 14
          parent:
            line_number: 14
            content: Marshal.restore(params[:oops])
          content: |
            Marshal.restore($<USER_INPUT>$<...>)
    - detector_id: ruby_lang_deserialization_of_user_input_user_input
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 1
          parent:
            line_number: 1
            content: params
          content: params
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 3
          parent:
            line_number: 3
            content: params
          content: params
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 5
          parent:
            line_number: 5
            content: params
          content: params
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 7
          parent:
            line_number: 7
            content: params
          content: params
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 9
          parent:
            line_number: 9
            content: params
          content: params
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 11
          parent:
            line_number: 11
            content: params
          content: params
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 13
          parent:
            line_number: 13
            content: params
          content: params
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_params.rb
          line_number: 14
          parent:
            line_number: 14
            content: params
          content: params
components: []


--

