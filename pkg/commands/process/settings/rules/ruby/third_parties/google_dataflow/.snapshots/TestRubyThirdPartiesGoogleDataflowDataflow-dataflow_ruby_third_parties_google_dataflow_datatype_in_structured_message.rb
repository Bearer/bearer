data_types:
    - name: IP address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_structured_message.rb
              line_number: 4
              field_name: ip_address
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_dataflow
      data_types:
        - name: IP address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_structured_message.rb
              line_number: 4
              parent:
                line_number: 4
                content: 'str_msg.message_text = "Current user: #{user.ip_address}"'
              field_name: ip_address
              object_name: user
              subject_name: User
    - detector_id: google_dataflow_message_text_classes
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_structured_message.rb
          line_number: 3
          parent:
            line_number: 3
            content: Google::Cloud::Dataflow::V1beta3::StructuredMessage.new
          content: |
            Google::Cloud::Dataflow::$<_>::StructuredMessage.new
components: []


--

