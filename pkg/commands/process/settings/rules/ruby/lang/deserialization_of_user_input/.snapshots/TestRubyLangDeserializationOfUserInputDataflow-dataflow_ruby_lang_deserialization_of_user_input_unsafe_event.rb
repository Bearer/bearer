risks:
    - detector_id: ruby_lang_deserialization_of_user_input
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
          line_number: 2
          parent:
            line_number: 2
            content: YAML.load(event["oops"])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
          line_number: 4
          parent:
            line_number: 4
            content: Psych.load(event["oops"])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
          line_number: 6
          parent:
            line_number: 6
            content: Syck.load(event["oops"])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
          line_number: 8
          parent:
            line_number: 8
            content: JSON.load(event["oops"])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
          line_number: 10
          parent:
            line_number: 10
            content: Oj.load(event["oops"])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
          line_number: 11
          parent:
            line_number: 11
            content: |-
                Oj.object_load(event["oops"]) do |json|
                  end
          content: |
            Oj.object_load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
          line_number: 14
          parent:
            line_number: 14
            content: Marshal.load(event["oops"])
          content: |
            $<LIBRARY>.load($<USER_INPUT>$<...>)$<...>
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
          line_number: 15
          parent:
            line_number: 15
            content: Marshal.restore(event["oops"])
          content: |
            Marshal.restore($<USER_INPUT>$<...>)
    - detector_id: ruby_lang_deserialization_of_user_input_user_input
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/deserialization_of_user_input/testdata/unsafe_event.rb
          line_number: 1
          parent:
            line_number: 1
            content: event
          content: |
            def $<_>($<!>event:, context:)
            end
components: []


--

