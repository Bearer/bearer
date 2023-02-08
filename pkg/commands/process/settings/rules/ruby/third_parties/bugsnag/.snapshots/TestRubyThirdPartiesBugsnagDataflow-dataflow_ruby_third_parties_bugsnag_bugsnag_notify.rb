data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/bugsnag_notify.rb
              line_number: 2
              field_name: email
              object_name: current_user
              subject_name: User
    - name: Fullname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/bugsnag_notify.rb
              line_number: 13
              field_name: name
              object_name: current_user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_bugsnag
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/bugsnag_notify.rb
          line_number: 7
          parent:
            line_number: 7
            content: |-
                Bugsnag.notify(exception) do |event|
                  # Adjust the severity of this error
                  event.severity = "warning"

                  # Add customer information to this event
                  event.add_metadata(:account, {
                    user_name: current_user.name,
                    paying_customer: true
                  })
                end
          content: |
            Bugsnag.notify($<_>) do |$<EVENT:identifier>|
              $<EVENT>.add_metadata($<...>$<DATA_TYPE>$<...>)
            end
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/bugsnag/testdata/bugsnag_notify.rb
          line_number: 7
          parent:
            line_number: 7
            content: |-
                Bugsnag.notify(exception) do |event|
                  # Adjust the severity of this error
                  event.severity = "warning"

                  # Add customer information to this event
                  event.add_metadata(:account, {
                    user_name: current_user.name,
                    paying_customer: true
                  })
                end
          content: |
            Bugsnag.notify($<_>) do |$<EVENT:identifier>|
              $<EVENT>.add_metadata($<...>$<DATA_TYPE>$<...>)
            end
components: []


--

