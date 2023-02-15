risks:
    - detector_id: ruby_lang_deserialization_of_user_input
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 1
          parent:
            line_number: 1
            content: YAML.load(request.env[:oops])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 3
          parent:
            line_number: 3
            content: Psych.load(request.env[:oops])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 5
          parent:
            line_number: 5
            content: Syck.load(request.env[:oops])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 7
          parent:
            line_number: 7
            content: JSON.load(request.env[:oops])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 9
          parent:
            line_number: 9
            content: Oj.load(request.env[:oops])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 10
          parent:
            line_number: 10
            content: |-
                Oj.object_load(request.env[:oops]) do |json|
                end
          content: |
            Oj.object_load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 13
          parent:
            line_number: 13
            content: Marshal.load(request.env[:oops])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 14
          parent:
            line_number: 14
            content: Marshal.restore(request.env[:oops])
          content: |
            Marshal.restore($<USER_INPUT>$<...>)
    - detector_id: ruby_lang_deserialization_of_user_input_user_input
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 1
          parent:
            line_number: 1
            content: request
          content: request
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 3
          parent:
            line_number: 3
            content: request
          content: request
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 5
          parent:
            line_number: 5
            content: request
          content: request
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 7
          parent:
            line_number: 7
            content: request
          content: request
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 9
          parent:
            line_number: 9
            content: request
          content: request
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 10
          parent:
            line_number: 10
            content: request
          content: request
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 13
          parent:
            line_number: 13
            content: request
          content: request
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_request.rb
          line_number: 14
          parent:
            line_number: 14
            content: request
          content: request
components: []


--

