data_types:
    - name: IP address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
              line_number: 2
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
              line_number: 5
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
              line_number: 8
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
              line_number: 11
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
              line_number: 12
              field_name: ip_address
              object_name: customer
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_dataflow
      data_types:
        - name: IP address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'param.value = "ip: #{customer.ip_address}"'
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
              line_number: 5
              parent:
                line_number: 5
                content: 'custom_metadata.value = "ip: #{customer.ip_address}"'
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
              line_number: 8
              parent:
                line_number: 8
                content: 'param2.value = "ip: #{customer.ip_address}"'
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
              line_number: 11
              parent:
                line_number: 11
                content: 'template_metadata.description ="ip: #{customer.ip_address}"'
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
              line_number: 12
              parent:
                line_number: 12
                content: 'template_metadata.name ="ip: #{customer.ip_address}"'
              field_name: ip_address
              object_name: customer
              subject_name: User
    - detector_id: google_dataflow_custom_metadata_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
          line_number: 4
          parent:
            line_number: 4
            content: Google::Cloud::Dataflow::V1beta3::ParameterMetadata::CustomMetadataEntry.new()
          content: |
            Google::Cloud::Dataflow::$<_>::ParameterMetadata::CustomMetadataEntry.new()
    - detector_id: google_dataflow_parameters_entry_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
          line_number: 1
          parent:
            line_number: 1
            content: Google::Cloud::Dataflow::V1beta3::CreateJobFromTemplateRequest::ParametersEntry.new()
          content: |
            Google::Cloud::Dataflow::$<_>::CreateJobFromTemplateRequest::ParametersEntry.new()
    - detector_id: google_dataflow_structured_message_parameter_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
          line_number: 7
          parent:
            line_number: 7
            content: Google::Cloud::Dataflow::V1beta3::StructuredMessage::Parameter.new()
          content: |
            Google::Cloud::Dataflow::$<_>::StructuredMessage::Parameter.new()
    - detector_id: google_dataflow_template_metadata_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/temp.rb
          line_number: 10
          parent:
            line_number: 10
            content: Google::Cloud::Dataflow::V1beta3::TemplateMetadata.new()
          content: |
            Google::Cloud::Dataflow::$<_>::TemplateMetadata.new()
components: []


--

