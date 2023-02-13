data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_config.rb
              line_number: 8
              field_name: email
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_config.rb
              line_number: 13
              field_name: email
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_config.rb
              line_number: 14
              field_name: email
              object_name: current_user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_dataflow
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_config.rb
              line_number: 8
              parent:
                line_number: 8
                content: 'config.metadata = { current_user_id: current_user.email }'
              field_name: email
              object_name: current_user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_config.rb
              line_number: 14
              parent:
                line_number: 14
                content: 'client_config.metadata = { current_user_id: current_user.email }'
              field_name: email
              object_name: current_user
              subject_name: User
    - detector_id: google_dataflow_client_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_config.rb
          line_number: 2
          parent:
            line_number: 2
            content: Google::Cloud::Dataflow.messages
          content: |
            Google::Cloud::Dataflow.$<METHOD>
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_config.rb
          line_number: 3
          parent:
            line_number: 3
            content: Google::Cloud::Dataflow::V1beta3::Messages::Client.new
          content: |
            Google::Cloud::Dataflow::$<_>::$<_>::Client.new
    - detector_id: google_dataflow_config
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_config.rb
          line_number: 6
          parent:
            line_number: 6
            content: client.configure
          content: |
            $<VAR>.configure
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_config.rb
          line_number: 7
          parent:
            line_number: 7
            content: config
          content: |
            $<VAR>.configure { |$<!>$<_:identifier>| }
components: []


--

